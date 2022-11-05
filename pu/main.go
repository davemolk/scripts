package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
)

type config struct {
	domains bool
	file string
	keys bool
	kv bool
	paths bool
	user bool
	values bool
	verbose bool
}

type pu struct {
	config config
}

func main() {
	var config config
	flag.BoolVar(&config.domains, "domains", false, "output domains")
	flag.StringVar(&config.file, "file", "", "name of file containing urls to parse")
	flag.BoolVar(&config.keys, "keys", false, "output keys")
	flag.BoolVar(&config.kv, "kv", false, "output keys and values")
	flag.BoolVar(&config.paths, "paths", false, "output paths")
	flag.BoolVar(&config.user, "user", false, "output username and password")
	flag.BoolVar(&config.values, "values", false, "output values")
	flag.BoolVar(&config.verbose, "v", false, "verbose output")
	flag.Parse()

	p := &pu{
		config: config,
	}
	
	if !p.inputValid(config) {
		log.Fatal("must choose something to parse")
	}
	
	ch, err := p.read(config)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case config.domains:
		for u := range p.domains(p.parsed(ch)) {
			fmt.Println(u)
		}
	case config.kv:
		for u := range p.kvMap(p.parsed(ch)) {
			fmt.Println(u)
		}
	
	case config.paths:
		for u := range p.paths(p.parsed(ch)) {
			fmt.Println(u)
		}
	case config.user:
		for u := range p.user(p.parsed(ch)) {
			fmt.Println(u)
		}
	}
}

func (p *pu) inputValid(config config) bool {
	if !config.domains && !config.keys && !config.kv && !config.paths && !config.user && !config.values {
		return false
	}
	return true
}

func (p *pu) domains(urls <-chan *url.URL) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for u := range urls {
			ch <- u.Host
		}
	}()
	return ch
}

func (p *pu) kvMap(urls <-chan *url.URL) <-chan map[string][]string {
	ch := make(chan map[string][]string)

	go func() {
		defer close(ch)
		for u := range urls {
			m, err := url.ParseQuery(u.RawQuery)
			if err != nil {
				if p.config.verbose{
					log.Printf("param parsing error: %v\n", err)
				}
				continue
			}
			ch <- m
		}
	}()
	return ch
}

// func getKeys(kv <-chan map[string][]string) <-chan []string {
// 	ch := make(chan []string)

// 	go func() {
// 		defer close(ch)
// 		keys := make([]string, len(kv))
// 		i := 0
// 		for k := range kv {
// 			for key := range k {

// 			}
// 		}
// 	}()
// }

func (p *pu) paths(urls <-chan *url.URL) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		for u := range urls {
			if u.Path != "" {
				ch <- u.Path
			}
		}
	}()
	return ch
}

func (p *pu) parsed(urls <-chan string) <-chan *url.URL {
	ch := make(chan *url.URL)

	go func() {
		defer close(ch)
		for u := range urls {
			s, err := url.ParseRequestURI(u)
			if err != nil {
				if p.config.verbose {
					log.Printf("parsing error: %v\n", err)
				}
				continue
			}
			ch <- s
		}
	}()

	return ch
}

func (p *pu) user(urls <-chan *url.URL) <-chan *url.Userinfo {
	ch := make(chan *url.Userinfo)

	go func() {
		defer close(ch)
		for u := range urls {
			if u.User != nil {
				ch <- u.User
			}
		}
	}()
	return ch
}

func (p *pu) read(config config) (<-chan string, error) {
	ch := make(chan string)
	var scanner *bufio.Scanner

	switch {
	case config.file != "":
		f, err := os.Open(config.file)
		if err != nil {
			return nil, fmt.Errorf("read error: %v", err)
		}
		scanner = bufio.NewScanner(f)
	default:
		scanner = bufio.NewScanner(os.Stdin)
	}
	
	go func() {
		defer close(ch)
		for scanner.Scan() {
			ch <- scanner.Text()
		}
		if err := scanner.Err(); err != nil && p.config.verbose {
			log.Printf("error for %s: %v\n", config.file, err)
		}
	}()
	
	return ch, nil
}
