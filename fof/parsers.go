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

func (f *fof) parseBing(data, control string) {
	sr := search{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}

	test := doc.Find("head title").Text()
	if test == control {
		fmt.Println("match")
	} else {
		fmt.Println("no match")
	}

	doc.Find("li.b_algo").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Find("h2 a").Attr("href"); ok {
			sr.URL = link
		} else {
			sr.URL = ""
		}
		blurb := s.Find("div.b_caption p").Text()
		sr.Blurb = blurb
		f.searches.store("music", sr)
	})
}

// eventually just pass in selectors to one function...
func (f *fof) parseGoogle(data, control string) {
	sr := search{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}

	test := doc.Find("head title").Text()
	if test == control {
		fmt.Println("match")
	} else {
		fmt.Println("no match")
	}

	doc.Find("div.g").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Find("a").Attr("href"); ok {
			sr.URL = link
		}
		blurb := s.Find("div[style='-webkit-line-clamp:2'] span").Text()
		sr.Blurb = blurb
		f.searches.store("music", sr)
	})
}

func (f *fof) parseYahoo(data, control string) {
	sr := search{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		fmt.Println(err)
	}

	test := doc.Find("head title").Text()
	if test == control {
		fmt.Println("match")
	} else {
		fmt.Println("no match")
	}

	doc.Find("div.algo").Each(func(i int, s *goquery.Selection) {
		if link, ok := s.Find("h3 > a").Attr("href"); ok {
			sr.URL = link
		}
		blurb := s.Find("div.compText").Text()
		sr.Blurb = blurb
		f.searches.store("music", sr)
	})
}
