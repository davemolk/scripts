package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	var dir string
	flag.StringVar(&dir, "", "", "directory to read files from")

	files, err := readDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var allChan []<-chan string
	for _, file := range files {
		ch, err := readFile(file)
		if err != nil {
			fmt.Println(err)
			continue
		}
		allChan = append(allChan, ch)
	}

	exit := make(chan struct{})
	chM := merge(allChan)

	go func() {
		defer close(exit)
		f, err := os.Create("results.txt")
		if err != nil {
			fmt.Println(err)
			return	
		}
		defer f.Close()
		for v := range chM {
			fmt.Fprintf(f, "%s\n", v)
		}
	}()

	<-exit
	fmt.Println("all done")
}

func merge(cs []<-chan string) <-chan string {
	var wg sync.WaitGroup
	out := make(chan string)

	send := func(c <-chan string) {
		defer wg.Done()
		for n := range c {
			fmt.Println(n)
			out <- n
		}
	}

	wg.Add(len(cs))

	for _, c := range cs {
		go send(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func readDir(dir string) ([]string, error) {
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			log.Fatal(err)
		}
		dir = wd
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var names []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".txt") {
			names = append(names, file.Name())
		}
	}
	return names, nil
}

func readFile(file string) (<-chan string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	ch := make(chan string)
	scanner := bufio.NewScanner(f)

	go func() {
		for scanner.Scan() {
			ch <- scanner.Text()
			if err := scanner.Err(); err != nil {
				fmt.Printf("error for %s: %v\n", file, err)
				break
			}

		}
		close(ch)
	}()

	return ch, nil
}