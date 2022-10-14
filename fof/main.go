package main

import (
	"flag"
	"fmt"
)

type config struct {
	query   string
	term    string
	file    string
	timeout int
}

type fof struct {
	config   config
	searches *searchesMap
	terms    []string
}

func main() {
	var config config
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

	qd := f.makeQueryData()

	// conditional to determine term vs terms

	urlB, controlB := f.makeQueryString(qd[1])

	fmt.Println(urlB, controlB)

	for i, u := range urlB {
		s, err := f.makeRequest(u, config.timeout)
		if err != nil {
			fmt.Println(err)
		}

		f.parseYahoo(s, controlB[i])
	}

	// s, err := f.makeRequest(urlB[0], config.timeout)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// f.parseYahoo(s, controlB[0])

	fmt.Println(f.searches.searches)
}
