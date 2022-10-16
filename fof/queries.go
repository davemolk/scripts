package main

import (
	"fmt"
	"strings"
	"sync"
)

type queryData struct {
	base   string
	spacer string
}

type parseData struct {
	blurbSelector string
	itemSelector  string
	linkSelector  string
}

func (f *fof) makeQueryData() []*queryData {
	var qd []*queryData

	bing := &queryData{
		base:   "https://bing.com/search?q=",
		spacer: "%20",
	}
	qd = append(qd, bing)

	// blocks a lot
	google := &queryData{
		base:   "https://www.google.com/search?q=",
		spacer: "+",
	}

	_ = google

	yahoo := &queryData{
		base:   "https://search.yahoo.com/search?p=",
		spacer: "+",
	}

	qd = append(qd, yahoo)

	return qd
}

func (f *fof) makeParseData() []*parseData {
	var pd []*parseData

	bing := &parseData{
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
	}
	pd = append(pd, bing)

	google := &parseData{
		blurbSelector: "div[style='-webkit-line-clamp:2'] span",
		itemSelector:  "div.g",
		linkSelector:  "a",
	}

	_ = google

	yahoo := &parseData{
		blurbSelector: "div.compText",
		itemSelector:  "div.algo",
		linkSelector:  "h3 > a",
	}

	pd = append(pd, yahoo)

	return pd
}

func (f *fof) makeQueryString(wg *sync.WaitGroup, data *queryData, term string, ch chan string) {
	defer wg.Done()
	cleanQ := strings.Replace(f.config.query, " ", data.spacer, -1)
	url := fmt.Sprintf("%s%s%s%s", data.base, cleanQ, data.spacer, term)
	// jenky, lol
	url = fmt.Sprintf("%sGETTERM%s", url, term)
	ch <- url
}

func (f *fof) makeSearchURLs(qdSlice [] *queryData) [2]chan string {
	var chans [2]chan string
	for i := range chans {
		chans[i] = make(chan string, len(f.terms))
	}

	var wg sync.WaitGroup
	for _, term := range f.terms {
		for i, qd := range qdSlice {
			wg.Add(1)
			go f.makeQueryString(&wg, qd, term, chans[i])
		}
	}

	wg.Wait()
	for i := range chans {
		close(chans[i])
	}

	return chans
}

func (f *fof) getAndParseData(pdSlice []*parseData, chans [2]chan string) {
	var wg sync.WaitGroup
	tokens := make(chan struct{}, f.config.concurrency)
	for i, ch := range chans {
		for u := range ch {
			wg.Add(1)
			tokens <- struct{}{}
			go func(i int, u string) {
				defer wg.Done()
				urlTerm := strings.Split(u, "GETTERM")
				s, err := f.makeRequest(urlTerm[0], f.config.timeout)
				if err != nil {
					fmt.Println("error")
					<-tokens
					return
				}
				<-tokens
				f.parseSearchResults(s, urlTerm[1], pdSlice[i])
			}(i, u)
		}
	}

	wg.Wait()
}