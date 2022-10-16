package main

import (
	"fmt"
	"log"
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
	log.Printf("parsing %s for %q", pd.name, term)
	sr := search{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Printf("goquery error for %s: %v\n", pd.name, err)
		return
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
		} else {
			log.Printf("unable to get link for %s\n", pd.name)
		}
		blurb := s.Find(pd.blurbSelector).Text()
		if blurb == "" {
			log.Printf("unable to get blurb for %s\n", pd.name) // check that it is this w/ no result
		}
		sr.Blurb = blurb
		f.searches.store(term, sr)
	})
}
