package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

type config struct {
	gophers int
	json    bool
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
	flag.IntVar(&config.timeout, "t", 5000, "timeout for requests (in ms, default 5000)")
	flag.BoolVar(&config.txt, "txt", false, "output results as txt (default false)")
	flag.StringVar(&config.url, "u", "", "url to get")
	flag.Parse()

	t := &tas{
		config: config,
	}

	if config.url == "" {
		log.Fatal("url cannot be empty")
	}

	t.client = t.makeClient()
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
}

func (t *tas) writeData(title string, data map[int][]string) {
	f, err := os.Create(title)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if t.config.txt {
		for i, d := range data {
			line := fmt.Sprintf("%d: %s", i, d)
			fmt.Fprintln(f, line)
		}
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("marshal error within writeData: %v", err)
		}
		_, err = f.Write(b)
		if err != nil {
			log.Fatalf("write error within writeData: %v", err)
		}
		err = f.Sync()
		if err != nil {
			log.Fatalf("sync error: %v", err)
		}
	}
}
