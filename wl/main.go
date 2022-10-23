package main

import (
	"flag"
	"fmt"
	"log"
	"regexp"
	"sync"
)

type config struct {
	timeout int
	url string
}

type wl struct {
	config
	noBlank *regexp.Regexp
}

func main() {
	var config config
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms)")
	flag.StringVar(&config.url, "u", "", "url to search")
	flag.Parse()

	noBlank := regexp.MustCompile(`\s{2,}`)
	w := &wl{
		config: config,
		noBlank: noBlank,
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
			fmt.Println(word)
		}(word)
	}
	wg.Wait()
}