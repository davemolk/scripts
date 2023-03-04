package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/davemolk/scripts/big"
)


func main() {
	var dir string
	var size int64
	flag.StringVar(&dir, "d", ".", "directory to start walk")
	flag.Int64Var(&size, "s", 0, "minimum file size")
	flag.Parse()

	c, err := big.Files(os.DirFS(dir), size)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(c)
}
