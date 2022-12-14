package main

import (
	"net/url"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

// searchesMap has the search term(s) as the key(s) and a
// nested map as the value(s). The nested map is in the
// form URL: blurb.
type searchesMap struct {
	mu       sync.Mutex
	searches map[string]map[string]string
}

func newSearchMap() *searchesMap {
	return &searchesMap{
		searches: make(map[string]map[string]string),
	}
}

func (s *searchesMap) store(term, url, blurb string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searches[term][url]; ok {
		return
	}
	s.searches[term][url] = blurb
}

func (f *fof) getAndParseData(pdSlice []*parseData, chans [6]chan string) {
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
					f.errorLog.Printf("error in makeRequest: %v\n", err)
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

func (f *fof) parseSearchResults(data, term string, pd *parseData) {
	f.infoLog.Printf("Parsing %s for %q", pd.name, term)
	localResults := make(map[string]string)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(data))
	if err != nil {
		f.errorLog.Printf("goquery error for %s: %v\n", pd.name, err)
		return
	}

	doc.Find(pd.itemSelector).Each(func(i int, s *goquery.Selection) {
		// no link, no point in getting blurb
		if link, ok := s.Find(pd.linkSelector).Attr("href"); !ok {
			f.errorLog.Printf("unable to get link for %s\n", pd.name)
			return
		} else {
			cleanLink := f.cleanLinks(link)
			blurb := s.Find(pd.blurbSelector).Text()
			if blurb == "" {
				f.errorLog.Printf("unable to get blurb for %s\n", pd.name)
			}
			cleanBlurb := f.cleanBlurb(blurb)
			localResults[cleanLink] = cleanBlurb
			f.searches.store(term, cleanLink, cleanBlurb)
		}

	})
}

func (f *fof) cleanBlurb(s string) string {
	cleanB := f.noBlank.ReplaceAllString(s, " ")
	cleanB = strings.TrimSpace(cleanB)
	cleanB = strings.ReplaceAll(cleanB, "\n", "")
	return cleanB
}

func (f *fof) cleanLinks(s string) string {
	u, err := url.QueryUnescape(s)
	if err != nil {
		f.errorLog.Printf("unable to clean %s: %v\n", s, err)
		return s
	}
	if strings.HasPrefix(u, "//duck") {
		removePrefix := strings.Split(u, "=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "&rut")
		u = removeSuffix[0]
	}
	if strings.HasPrefix(u, "https://r.search.yahoo.com/") {
		removePrefix := strings.Split(u, "/RU=")
		u = removePrefix[1]
		removeSuffix := strings.Split(u, "/RK=")
		u = removeSuffix[0]
	}
	return u
}
