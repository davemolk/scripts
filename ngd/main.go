package main

import (
	"flag"
)

type config struct {
	exclude string
	filetype string
	inTitle bool
	inURL bool
	query string
	queryExact string
	site string
	timeout int
}

type ngd struct {
	config config
}

func main() {
	var config config
	flag.StringVar(&config.filetype, "ft", "", "file type")
	flag.BoolVar(&config.inURL, "url", false, "term to find in URL")
	flag.StringVar(&config.query, "q", "", "search query")
	flag.StringVar(&config.queryExact, "qe", "", "search query (exact matching)")
	flag.StringVar(&config.site, "site", "", "site/domain to search")
	flag.StringVar(&config.exclude, "no", "", "site/domain to exclude")
	flag.IntVar(&config.timeout, "t", 5000, "timeout for request")
	flag.BoolVar(&config.inTitle, "title", false, "term to find in site title")
	flag.Parse()

	n := &ngd{
		config: config,
	}

	n.validateInput(config)

	qd := n.makeQueryData()
	pd := n.makeParseData()
	url := n.makeQueryString(qd)
	
	n.getAndParseData(url, config.timeout, pd)
}