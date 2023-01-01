package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

func makePaths() []string {
	paths, err := input()
	if err != nil {
		log.Fatal(err)
	}
	for _, p := range paths {
		paths = append(paths, addPairs(p)...)
		paths = append(paths, addLeadings(p))
		paths = append(paths, addTrailings(p)...)
	}
	
	return paths
}

func createRequests(url string, paths []string) {
	headers := headers()

	ch := make(chan *http.Request, len(headers) * len(paths))
	var wg sync.WaitGroup
	for _, h := range headers {
		for _, p := range paths {
			wg.Add(1)
			go func(h []string, p string) {
				defer wg.Done()
				req, err := http.NewRequest(http.MethodGet, url + p, nil)
				if err != nil {
					fmt.Println(err)
					return
				}
				req.Header.Set(h[0], h[1])
				ch <- req
			}(h, p)
		}
	}
	wg.Wait()
	close(ch)
	for c := range ch {
		fmt.Println(c)
	}
}

func headers() [][]string {
	headers := []string{
		"X-Custom-IP-Authorization", "X-Forwarded-For", 
		"X-Forward-For", "X-Remote-IP", "X-Originating-IP", 
		"X-Remote-Addr", "X-Client-IP", "X-Real-IP",
	}

	values := []string{
		"localhost", "localhost:80", "localhost:443", 
		"127.0.0.1", "127.0.0.1:80", "127.0.0.1:443", 
		"2130706433", "0x7F000001", "0177.0000.0000.0001", 
		"0", "127.1", "10.0.0.0", "10.0.0.1", "172.16.0.0", 
		"172.16.0.1", "192.168.1.0", "192.168.1.1",
	}

	var header [][]string
	for _, h := range headers {
		for _, v := range values {
			header = append(header, []string{h, v})
		}
	}

	return header
}

func addLeadings(u string) string {
	return fmt.Sprintf("%v/%s", "/%2e", u)
}

func addTrailings(u string) []string {
	trail := []string{
		"/", "..;/", "/..;/", "%20", "%09", "%00", 
		".json", ".css", ".html", "?", "??", "???", 
		"?testparam", "#", "#test", "/.",
	}
	var urls []string
	for _, t := range trail {
		urls = append(urls, fmt.Sprintf("/%s%v", u, t))
	}
	return urls
}

func addPairs(u string) []string {
	var urls []string

	urls = append(urls, fmt.Sprintf("//%s//", u))
	urls = append(urls, fmt.Sprintf("/./%s/./", u))
	return urls
}