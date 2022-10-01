package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

// job stores a port number and any error encountered when
// scanning that port.
type job struct {
	Port int   `json:"port"`
	Err  error `json:"err,omitempty"`
}

// worker checks a port number and sends the results of the scan
// to the results channel.
func worker(addr string, timeout int, ports <-chan job, results chan<- job) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", addr, p.Port)
		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			p.Err = err
			results <- p
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	var addr string
	var errorReport bool
	var gophers int
	var ports int
	var timeout int

	flag.StringVar(&addr, "a", "", "address to scan")
	flag.BoolVar(&errorReport, "e", false, "report any errors")
	flag.IntVar(&gophers, "g", 20, "number of goroutines to use")
	flag.IntVar(&ports, "p", 1024, "upper boundary of ports to scan")
	flag.IntVar(&timeout, "t", 2000, "timeout (in ms)")

	flag.Parse()

	if addr == "" {
		log.Fatal("no address specified")
	}

	errorLog := log.New(os.Stderr, "ERROR\t", log.Ltime)

	jobs := make(chan job, ports)
	results := make(chan job, ports)

	for i := 1; i <= gophers; i++ {
		go worker(addr, timeout, jobs, results)
	}

	for j := 1; j <= ports; j++ {
		jobs <- job{Port: j}
	}

	close(jobs)

	var openPorts []job
	var errorPorts []job

	for p := 1; p <= ports; p++ {
		port := <-results
		if port.Err == nil {
			openPorts = append(openPorts, port)
		} else if errorReport {
			errorPorts = append(errorPorts, port)
		}
	}

	if errorReport {
		err := output("errorResults.json", errorPorts)
		if err != nil {
			errorLog.Println(err)
		}
	}

	for _, port := range openPorts {
		log.Printf("%d open\n", port.Port)
	}

	err := output("results.json", openPorts)
	if err != nil {
		errorLog.Println("output error:", err)
	}
}

// output takes a string for a file name and a job slice,
// marshals the job slice, and passes the data to writeData.
func output(name string, data []job) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return writeData(name, b)
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
