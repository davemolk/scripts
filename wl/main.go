package main

import (
	"flag"
	"log"
	"regexp"
	"sync"
)

type config struct {
	filter string
	minCount  int
	minLength int
	timeout   int
	url       string
}

type wl struct {
	config
	noBlank *regexp.Regexp
	wordMap *WordMap
}

func main() {
	var config config
	flag.StringVar(&config.filter, "f", "", "file name containing words to filter out of results")
	flag.IntVar(&config.minCount, "c", 0, "minimum count to include word in results")
	flag.IntVar(&config.minLength, "len", 0, "minimum word length to consider")
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms)")
	flag.StringVar(&config.url, "u", "", "url to search")
	flag.Parse()

	noBlank := regexp.MustCompile(`\s{2,}`)
	wordMap := newWordMap()

	w := &wl{
		config:  config,
		noBlank: noBlank,
		wordMap: wordMap,
	}

	g, err := w.makeRequest(config.url, config.timeout)
	if err != nil {
		log.Fatal(err)
	}

	words := w.processData(g)
	var wg sync.WaitGroup
	for _, word := range words {
		wg.Add(1)
		go func(word string) {
			defer wg.Done()
			word = w.removePunctuation(word)
			if len(word) >= config.minLength {
				w.wordMap.add(word)
			}
		}(word)
	}
	wg.Wait()
	
	if config.filter != "" {
		w.filterTerms()
	}
	
	keys := w.wordMap.sort()
	keysCount := w.dropLowCount(keys)

	w.wordMap.write(keysCount, "wl/results.txt")
}
