package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (n *ngd) getAndParseData(url string, timeout int, pd *parseData) {
	log.Println("getting", url)

	s, err := n.makeRequest(url, timeout)
	if err != nil {
		log.Println(err)
	}
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(s))
	if err != nil {
		log.Fatal("failed to parse html")
	}

	doc.Find(pd.itemSelector).Each(func(i int, s *goquery.Selection) {
		// no link, no point in getting blurb
		if link, ok := s.Find(pd.linkSelector).Attr("href"); !ok {
			return
		} else {
			fmt.Println(link)
			blurb := s.Find(pd.blurbSelector).Text()
			fmt.Println(blurb)
			fmt.Println()
		}
	})
}