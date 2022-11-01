package main

import (
	"fmt"
	"strings"
)

type queryData struct {
	base string
	colon string
	spacer string
}

type parseData struct {
	blurbSelector string
	itemSelector  string
	linkSelector  string
	name          string
}

func (n *ngd) makeQueryData() *queryData {
	ddg := &queryData{
		base: "https://html.duckduckgo.com/html?q=",
		colon: "%3A",
		spacer: "+", // %20 on fof...
	}

	return ddg
}

func (n *ngd) makeQueryString(qd *queryData) string {
	var components []string
	var cleanedQuery string
	switch {
	case n.config.query != "":
		cleanedQuery = strings.Replace(n.config.query, " ", qd.spacer, -1)
		components = append(components, cleanedQuery)
	case n.config.queryExact != "":
		cleanedQuery = strings.Replace(n.config.queryExact, " ", qd.spacer, -1)
		cleanedQuery = fmt.Sprintf("\"%s\"", cleanedQuery)
		components = append(components, cleanedQuery)
	}

	if n.config.inTitle {
		intitle := fmt.Sprintf("intitle%s%s", qd.colon, cleanedQuery)
		components = append(components, intitle)
	}
	
	if n.config.inURL {
		inurl := fmt.Sprintf("inurl%s%s", qd.colon, cleanedQuery)
		components = append(components, inurl)
	}

	if n.config.filetype != "" {
		filetype := fmt.Sprintf("filetype%s%s", qd.colon, n.config.filetype)
		components = append(components, filetype)
	}

	if n.config.site != "" {
		site := fmt.Sprintf("site%s%s", qd.colon, n.config.site)
		components = append(components, site)
	}

	if n.config.exclude != "" {
		exclude := fmt.Sprintf("-site%s%s", qd.colon, n.config.exclude)
		components = append(components, exclude)
	}

	params := strings.Join(components, "+")
	return fmt.Sprintf("%s%s", qd.base, params)
}

func (n *ngd) makeParseData() *parseData {
	duck := &parseData{
		blurbSelector: "div.links_main > a",
		itemSelector:  "div.web-result",
		linkSelector:  "div.links_main > a",
		name:          "duck",
	}
	return duck
}

