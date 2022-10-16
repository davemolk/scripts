package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// readInput takes in the file name for a list of terms and returns
// a string slice containing those terms.
func (f *fof) readInput(name string) ([]string, error) {
	var terms []string
	n, err := os.Open(name)
	if err != nil {
		return terms, err
	}
	defer n.Close()

	scanner := bufio.NewScanner(n)
	for scanner.Scan() {
		terms = append(terms, scanner.Text())
	}
	return terms, scanner.Err()
}

// getTerms looks at the user flag input, determines whether a single
// term or a file name for a list of terms has been selected, and 
// adds the appropriate field to the fof struct instance.
func (f *fof) getTerms() {
	var terms []string
	switch {
	case f.config.file != "":
		terms, err := f.readInput(f.config.file)
		if err != nil {
			log.Fatalf("unable to get terms: %v",err)
		}
		f.terms = terms
	case f.config.term != "":
		terms = append(terms, f.config.term)
		f.terms = terms
	default:
		fmt.Println("no search terms supplied")
	}
}
