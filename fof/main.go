package main

import (
	"flag"
	"fmt"
)

type config struct {
	concurrency int
	query       string
	term        string
	file        string
	timeout     int
}

type fof struct {
	config   config
	searches *searchesMap
	terms    []string
}

func main() {
	var config config
	flag.IntVar(&config.concurrency, "c", 10, "max number of goroutines to use at any given time")
	flag.StringVar(&config.query, "q", "", "search target (please enclose phrases in quotes)")
	flag.StringVar(&config.term, "term", "", "term to search for")
	flag.StringVar(&config.file, "file", "", "file name containing a list of terms")
	flag.IntVar(&config.timeout, "t", 5000, "timeout (in ms, default 5000)")

	flag.Parse()

	searches := newSearchMap()

	f := &fof{
		config:   config,
		searches: searches,
	}

	f.getTerms()
	qdSlice := f.makeQueryData()
	pdSlice := f.makeParseData()
	// length will be however many search engines we're using
	// var chans [2]chan string
	// for i := range chans {
	// 	chans[i] = make(chan string, len(f.terms))
	// }

	// var wg sync.WaitGroup
	// for _, term := range f.terms {
	// 	for i, q := range qd {
	// 		wg.Add(1)
	// 		go f.makeQueryString(&wg, q, term, chans[i])
	// 	}
	// }

	// wg.Wait()
	// for i := range chans {
	// 	close(chans[i])
	// }

	chans := f.makeSearchURLs(qdSlice)

	f.getAndParseData(pdSlice, chans)

	// tokens := make(chan struct{}, config.concurrency)
	// for i, ch := range chans {
	// 	for u := range ch {
	// 		wg.Add(1)
	// 		tokens <- struct{}{}
	// 		go func(i int, u string) {
	// 			defer wg.Done()
	// 			urlTerm := strings.Split(u, "GETTERM")
	// 			s, err := f.makeRequest(urlTerm[0], config.timeout)
	// 			if err != nil {
	// 				fmt.Println("error")
	// 				<-tokens
	// 				return
	// 			}
	// 			<-tokens
	// 			f.parseSearchResults(s, urlTerm[1], pd[i])
	// 		}(i, u)
	// 	}
	// }

	// wg.Wait()

	fmt.Println(f.searches.searches["music"])
}
