package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

func main() {
	var dir string
	var ext string
	flag.StringVar(&dir, "dir", "", "directory to read files from")
	flag.StringVar(&ext, "ext", "txt", "file extension (csv or txt")
	flag.Parse()

	files, err := readDir(dir, ext)
	if err != nil {
		log.Fatal(err)
	}

	var inputChans []<-chan []string
	for _, file := range files {
		ch, err := readFile(file, ext)
		if err != nil {
			fmt.Println(err)
			continue
		}
		inputChans = append(inputChans, ch)
	}

	exit := make(chan struct{})
	mergeCh := merge(inputChans)

	go func() {
		defer close(exit)
		switch {
		case ext == "txt":
			writeTxt(mergeCh)
		case ext == "csv":
			writeCSV(mergeCh)
		}
	}()

	<-exit
	fmt.Println("all done!")
}

func writeTxt(mergeCh <-chan []string) {
	f, err := os.Create("results.txt")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	for v := range mergeCh {
		fmt.Fprintf(f, "%s\n", v[0])
	}
}

func writeCSV(mergeCh <-chan []string) {
	f, err := os.Create("results.csv")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	var data [][]string
	for v := range mergeCh {
		row := v
		data = append(data, row)
	}
	w.WriteAll(data)
}

func merge(cs []<-chan []string) <-chan []string {
	var wg sync.WaitGroup
	wg.Add(len(cs))

	out := make(chan []string)

	send := func(c <-chan []string) {
		defer wg.Done()
		for n := range c {
			out <- n
		}
	}

	for _, c := range cs {
		go send(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func readDir(dir, ext string) ([]string, error) {
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
		if strings.HasSuffix(file.Name(), ext) {
			names = append(names, file.Name())
		}
	}
	return names, nil
}

func readFile(file, ext string) (<-chan []string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("read error: %w", err)
	}

	ch := make(chan []string)

	switch {
	case ext == "csv":
		c := csv.NewReader(f)
		go func() {
			defer close(ch)
			for {
				record, err := c.Read()
				if err != nil {
					if errors.Is(err, io.EOF) {
						return
					}
					log.Fatal(err)
				}
				ch <- record
			}
		}()
	case ext == "txt":
		scanner := bufio.NewScanner(f)
		go func() {
			defer close(ch)
			for scanner.Scan() {
				var lines []string
				lines = append(lines, scanner.Text())
				ch <- lines
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("error for %s: %v\n", file, err)
			}
		}()
	}

	return ch, nil
}
