package main

import (
	"flag"
	"log"
	"os"
)

type config struct {
	concurrency  int
	searchTarget string
	term         string
	file         string
	timeout      int
}

type fof struct {
	config   config
	errorLog *log.Logger
	infoLog  *log.Logger
	searches *searchesMap
	terms    []string
}

func main() {
	var config config
	flag.IntVar(&config.concurrency, "c", 10, "max number of goroutines to use at any given time")
	flag.StringVar(&config.searchTarget, "s", "", "search target (please enclose phrases in quotes)")
	flag.StringVar(&config.term, "term", "", "term to search for")
	flag.StringVar(&config.file, "file", "", "file name containing a list of terms")
	flag.IntVar(&config.timeout, "t", 5000, "timeout (in ms, default 5000)")

	flag.Parse()

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime|log.Lshortfile)
	infoLog := log.New(os.Stdout, "INFO\t", log.Ltime)

	searches := newSearchMap()

	f := &fof{
		config:   config,
		errorLog: errorLog,
		infoLog:  infoLog,
		searches: searches,
	}

	f.getTerms()
	for _, t := range f.terms {
		searches.searches[t] = make(map[string]string)
	}

	qdSlice := f.makeQueryData()
	pdSlice := f.makeParseData()

	chans := f.makeSearchURLs(qdSlice)

	f.getAndParseData(pdSlice, chans)

}
