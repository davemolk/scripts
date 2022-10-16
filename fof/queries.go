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
	name string
}

func (f *fof) makeQueryData() []*queryData {
	var qdSlice []*queryData

	ask := &queryData{
		base:   "https://www.ask.com/web?q=",
		spacer: "%20",
	}
	qdSlice = append(qdSlice, ask)

	bing := &queryData{
		base:   "https://bing.com/search?q=",
		spacer: "%20",
	}
	qdSlice = append(qdSlice, bing)

	brave := &queryData{
		base:   "https://search.brave.com/search?q=",
		spacer: "+",
	}
	qdSlice = append(qdSlice, brave)

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
	qdSlice = append(qdSlice, yahoo)

	yandex := &queryData{
		base: "https://yandex.com/search/?text=",
		spacer: "+",
	}
	qdSlice = append(qdSlice, yandex)

	return qdSlice
}

func (f *fof) makeParseData() []*parseData {
	var pdSlice []*parseData

	ask := &parseData{
		blurbSelector: "div.PartialSearchResults-item p",
		itemSelector:  "div.PartialSearchResults-item",
		linkSelector:  "a.PartialSearchResults-item-title-link",
		name: "ask",
	}
	pdSlice = append(pdSlice, ask)

	bing := &parseData{
		blurbSelector: "div.b_caption p",
		itemSelector:  "li.b_algo",
		linkSelector:  "h2 a",
		name: "bing",
	}
	pdSlice = append(pdSlice, bing)

	brave := &parseData{
		blurbSelector: "div.snippet-content p.snippet-description",
		itemSelector:  "div.fdb",
		linkSelector:  "div.fdb > a.result-header",
		name: "brave",
	}
	pdSlice = append(pdSlice, brave)

	google := &parseData{
		blurbSelector: "div[style='-webkit-line-clamp:2'] span",
		itemSelector:  "div.g",
		linkSelector:  "a",
		name: "google",
	}
	_ = google

	yahoo := &parseData{
		blurbSelector: "div.compText",
		itemSelector:  "div.algo",
		linkSelector:  "h3 > a",
		name: "yahoo",
	}
	pdSlice = append(pdSlice, yahoo)

	yandex := &parseData{
		blurbSelector: "div.TextContainer",
		itemSelector: "li.serp-item",
		linkSelector: "div.VanillaReact > a",
		name: "yandex",
	}
	pdSlice = append(pdSlice, yandex)

	return pdSlice
}

func (f *fof) makeQueryString(wg *sync.WaitGroup, data *queryData, term string, ch chan string) {
	defer wg.Done()
	cleanQ := strings.Replace(f.config.query, " ", data.spacer, -1)
	url := fmt.Sprintf("%s%s%s%s", data.base, cleanQ, data.spacer, term)
	// jenky, lol
	url = fmt.Sprintf("%sGETTERM%s", url, term)
	ch <- url
}

func (f *fof) makeSearchURLs(qdSlice [] *queryData) [5]chan string {
	var chans [5]chan string
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

func (f *fof) getAndParseData(pdSlice []*parseData, chans [5]chan string) {
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
					fmt.Printf("error in makeRequest: %v\n", err)
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