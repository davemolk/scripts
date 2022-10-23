package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (w *wl) processData(doc *goquery.Document) []string {
	doc.Find("script").Remove()
	doc.Find("style").Remove()
	body := doc.Text()
	body = w.noBlank.ReplaceAllString(body, " ")
	body = strings.Replace(body, "\n", "", -1)
	
	return strings.Split(body, " ")
}

func (w *wl) removePunctuation(word string) string {
	punc := w.getPunctuation()
	for _, p := range punc {
		word = strings.TrimSuffix(word, p)
	}
	return strings.ToLower(word)
}

// prob add more
func (w *wl) getPunctuation() []string {
	return []string{
		",",
		".",
		":",
		";",
		"!",
		"?",
	}
}