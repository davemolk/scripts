package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"time"
)

// results holds a mutex and a slice of ints
// that stores any open ports found during the scan.
type results struct {
	mu   sync.Mutex
	Open []int `json:"open"`
}

// add creates a lock on on the results struct, appends
// the found port to the Open field, and releases the lock.
func (r *results) add(i int) {
	r.mu.Lock()
	r.Open = append(r.Open, i)
	r.mu.Unlock()
}

func main() {
	var addr string
	var gophers int
	var ports int
	var timeout int

	flag.StringVar(&addr, "a", "", "address to scan")
	flag.IntVar(&gophers, "g", 20, "number of gophers")
	flag.IntVar(&ports, "p", 1024, "upper limit of port numbers to scan")
	flag.IntVar(&timeout, "t", 2000, "timeout (in ms)")
	flag.Parse()

	if addr == "" {
		log.Fatal("no address specified")
	}

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime)

	r := &results{}

	var wg sync.WaitGroup
	tokens := make(chan struct{}, gophers)
	for i := 1; i <= ports; i++ {
		wg.Add(1)
		tokens <- struct{}{}
		go func(j int) {
			defer wg.Done()
			address := fmt.Sprintf("%s:%d", addr, j)
			conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Millisecond)
			if err != nil {
				<-tokens
				return
			}
			defer conn.Close()
			r.add(j)
			<-tokens
		}(i)
	}

	wg.Wait()

	b, err := json.Marshal(r.Open)
	if err != nil {
		errorLog.Fatalf("Marshal error: %v\n", err)
	}

	err = writeData("scanResults.json", b)
	if err != nil {
		errorLog.Fatalf("writeData error: %v\n", err)
	}
}

// writeData takes in a string for a file name and a byte slice
// and writes the data to a file. Any error in the process will be
// returned.
func writeData(name string, data []byte) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(data)
	if err != nil {
		return err
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}
