package main

import (
	"bufio"
	"log"
	"os"
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

func (w *wl) filterTerms() {
	terms, err := w.readInput(w.config.filter)
	if err != nil {
		log.Println(err)
	}
	for _, term := range terms {
		if _, ok := w.wordMap.words[term]; ok {
			w.wordMap.delete(term)
		}
	}
}

func (w *wl) readInput(name string) ([]string, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var terms []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		terms = append(terms, scanner.Text()) 
	}
	return terms, scanner.Err()
}