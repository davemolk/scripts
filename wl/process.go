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

// add more?
func (w *wl) getPunctuation() []string {
	return []string{
		",",
		".",
		":",
		";",
		"!",
		"?",
		"—",
	}
}

func (w *wl) dropLowCount(keys []string) []string {
	if w.config.minCount > 0 {
		for i, key := range keys {
			if w.wordMap.words[key] < w.config.minCount {
				keys = keys[:i]
				break
			}
		}
	}
	return keys
}
