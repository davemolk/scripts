package main

import (
	"bufio"
	"flag"
	"os"
	"strings"
)

func main() {
	var url string
	flag.StringVar(&url, "u", "https://www.example.com", "base url")
	flag.Parse()
	
	paths := makePaths()
	createRequests(url, paths)
	
}

func input() ([]string, error) {
	var lines []string
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		if s.Text() == "q" {
			break
		}
		lines = append(lines, strings.TrimPrefix(s.Text(), "/"))
	}
	return lines, s.Err()
}