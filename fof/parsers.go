package main

import (
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
	defer s.mu.Unlock()
	s.searches[term] = append(s.searches[term], search)
}

func (f *fof) parseSearchResults(data, term string, pd *parseData) {
	f.infoLog.Printf("Parsing %s for %q", pd.name, term)
	sr := search{}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		f.errorLog.Printf("goquery error for %s: %v\n", pd.name, err)
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
			f.errorLog.Printf("unable to get link for %s\n", pd.name)
		}
		blurb := s.Find(pd.blurbSelector).Text()
		if blurb == "" {
			f.errorLog.Printf("unable to get blurb for %s\n", pd.name) // check that it is this w/ no result
		}
		var cleaned string
		if pd.name == "brave" {
			cleaned = f.cleanBlurb(blurb, true)
		} else {
			cleaned = f.cleanBlurb(blurb, false)
		}
		sr.Blurb = cleaned
		f.searches.store(term, sr)
	})
}

func (f *fof) cleanBlurb(s string, b bool) string {
	cleaned := strings.TrimSpace(s)
	cleaned = strings.ReplaceAll(cleaned, "\n", "")
	if b {
		cleaned = strings.ReplaceAll(cleaned, "           ", "") // thanks Brave
	}
	return cleaned
}
