package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

// fileByStatusCode takes in a file name and a string slice
// and writes the data to a file.
func (t *tas) fileByStatusCode(name string, data []string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	for _, d := range data {
		fmt.Fprintln(f, d)
	}
	
}

// writeData takes in a file name and a map[int][]string and writes
// the data either to a JSON file or a txt file, depending on what was
// specified by the flags. 
func (t *tas) writeData(name string, data map[int][]string) {
	f, err := os.Create(name)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	if t.config.txt {
		for i, d := range data {
			line := fmt.Sprintf("%d: %s", i, d)
			fmt.Fprintln(f, line)
		}
	} else {
		b, err := json.Marshal(data)
		if err != nil {
			log.Fatalf("marshal error within writeData: %v", err)
		}
		_, err = f.Write(b)
		if err != nil {
			log.Fatalf("write error within writeData: %v", err)
		}
		err = f.Sync()
		if err != nil {
			log.Fatalf("sync error: %v", err)
		}
	}
}
