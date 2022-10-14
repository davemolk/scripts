package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func (f *fof) readInput(name string) ([]string, error) {
	var names []string
	n, err := os.Open(name)
	if err != nil {
		return names, err
	}
	defer n.Close()

	scanner := bufio.NewScanner(n)
	for scanner.Scan() {
		names = append(names, scanner.Text())
	}
	return names, scanner.Err()
}

func (f *fof) getTerms() {
	var terms []string
	switch {
	case f.config.file != "":
		terms, err := f.readInput(f.config.file)
		if err != nil {
			log.Fatal(err)
		}
		f.terms = terms
	case f.config.term != "":
		terms = append(terms, f.config.term)
		f.terms = terms
	default:
		fmt.Println("no search terms supplied")
	}
}
