package main

import (
	"fmt"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

type search struct {
	Blurb string
	URL   string
}

type searchesMap struct {
	mu       sync.Mutex
	searches map[string][]search
}

// some type of selector map thing

func newSearchMap() *searchesMap {
	return &searchesMap{
		searches: make(map[string][]search),
	}
}

func (s *searchesMap) store(term string, search search) {
	s.mu.Lock()
	s.searches[term] = append(s.searches[term], search)
	s.mu.Unlock()
}

func (f *fof) parseSearchResults(data, term string, pd *parseData) {
	sr := search{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}

	// test := doc.Find("head title").Text()
	// if test == control {
	// 	fmt.Println("match")
	// } else {
	// 	fmt.Println("no match")
	// }

	doc.Find(pd.itemSelector).Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Find(pd.linkSelector).Attr("href"); ok {
			sr.URL = link
		}
		blurb := s.Find(pd.blurbSelector).Text()
		sr.Blurb = blurb
		f.searches.store(term, sr)
	})
}
