package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

type config struct {
	gophers int
	json    bool
	redirect bool
	timeout int
	txt     bool
	url     string
}

type tas struct {
	client  *http.Client
	config  config
	results *statusMap
}

type statusMap struct {
	mu      sync.Mutex
	results map[int][]string
}

func newStatusMap() *statusMap {
	return &statusMap{
		results: make(map[int][]string),
	}
}

func (s *statusMap) add(code int, url string) {
	s.mu.Lock()
	s.results[code] = append(s.results[code], url)
	s.mu.Unlock()
}

func main() {
	var config config
	flag.IntVar(&config.gophers, "g", 10, "number of gophers")
	flag.BoolVar(&config.json, "json", true, "output results as json (default true)")
	flag.BoolVar(&config.redirect, "r", true, "allow redirects (default true)")
	flag.IntVar(&config.timeout, "t", 5000, "request timeout (in ms, default 5000)")
	flag.BoolVar(&config.txt, "txt", false, "output results as txt (default false)")
	flag.StringVar(&config.url, "u", "", "url to get")
	flag.Parse()

	start := time.Now()

	t := &tas{
		config: config,
	}

	if config.url == "" {
		log.Fatal("url cannot be empty")
	}

	t.client = t.makeClient(config.redirect)
	t.results = newStatusMap()

	s, err := t.getURLs(config.url, config.timeout)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	tokens := make(chan struct{}, config.gophers)

	// skip key/metadata
	for _, v := range s[1:] {
		tokens <- struct{}{}
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			t.testURLs(url)
			<-tokens
		}(v[0])
	}
	wg.Wait()

	var title string
	if config.txt {
		title = "tas/results.txt"
	} else {
		title = "tas/results.json"
	}

	t.writeData(title, t.results.results)
	
	for i, v := range t.results.results {
		wg.Add(1)
		go func (i int, v []string) {
			defer wg.Done()
			name := fmt.Sprintf("tas/%d.txt", i)
			t.fileByStatusCode(name, v)
		}(i, v)
	}

	wg.Wait()

	log.Printf("Took: %f seconds\n", time.Since(start).Seconds())
}
