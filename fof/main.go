package main

import (
	"flag"
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

	chans := f.makeSearchURLs(qdSlice)

	f.getAndParseData(pdSlice, chans)

	// fmt.Println(f.searches.searches)
}
