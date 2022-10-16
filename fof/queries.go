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

	// blocks a lot
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
