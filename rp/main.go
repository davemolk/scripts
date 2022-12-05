package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
)

func main() {
	var d bool
	var kv string
	flag.BoolVar(&d, "d", false, "decode query string(s)")
	flag.StringVar(&kv, "kv", "", "given key, replace value(s)")
	flag.Parse()

	if !d && kv == "" {
		log.Fatal("must select d or v")
	}
	
	exit := make(chan struct{})
	if d {
		go decode(exit)
	} else {
		go replaceValue(exit, kv)
	}

	<-exit
}

func decode(exit chan struct{}) {
	defer close(exit)
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		d, err := url.QueryUnescape(s.Text())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(d)
	}
	if err := s.Err(); err != nil {
		fmt.Println(err)
	}
}

func replaceValue(exit chan struct{}, kv string) {
	defer close(exit)

	p := strings.Split(kv, "=")
	if len(p) != 2 {
		log.Fatal("must supply input as key=value")
	}
	key := p[0]
	value := p[1]

	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		u, err := url.ParseRequestURI(s.Text())
		if err != nil {
			log.Println(err)
			continue
		}
		params := u.Query()

		params.Set(key, value)
		
		u.RawQuery = params.Encode()
		fmt.Println(u.String())
	}
	if err := s.Err(); err != nil {
		fmt.Println(err)
	}
}