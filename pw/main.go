package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	ch, err := read()
	if err != nil {
		log.Fatal(err)
	}
	var result []string
	for c := range specials(nums(lower(upper(ch)))) {
		result = append(result, c)
	}
	noSpace := strings.Replace(strings.Join(result, ""), " ", "", -1)
	fmt.Println(noSpace)
	// check at least 1 number, at least 1 lower, 1 upper, 1 special char, maybe weigh the 
	// things differently (instead of 50/50 chance of swapping?)
}

func read() (<-chan string, error) {
	ch := make(chan string)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(os.Stdin)
		str, err := reader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		for _, c := range strings.Split(str, "") {
			ch <- c
		}
	}()
	
	return ch, nil
}

func gimmeNum(num int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(num)
}

func upper(s <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for char := range s {
			if gimmeNum(2) == 1 {
				char = strings.ToUpper(char)
			}
			ch <- char
		}
	}()
	return ch
}

func lower(s <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for char := range s {
			if gimmeNum(2) == 1 {
				char = strings.ToLower(char)
			}
			ch <- char
		}
	}()
	return ch
}

func nums(s <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		nums := getNums()
		for char := range s {
			if gimmeNum(2) == 1 {
				char = nums[gimmeNum(len(nums))]
			}
			ch <- char
		}
	}()
	return ch
}

func specials(s <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		sp := getSpecials()
		for char := range s {
			if gimmeNum(2) == 1 {
				char = sp[gimmeNum(len(sp))]
			}
			ch <- char
		}
	}()
	return ch
}

func getSpecials() []string {
	return []string{
		`"`, " ", "!", "#", "$", "%", "&", "'", "(", ")", "*", "+", ",", "-", ".", "/", ":", ";", "<", "=", ">", "?", "@", "[", `\`, "`", "]", "^", "_", "{", "|", "}", "~",
	}
}

func getNums() []string {
	return []string{
		"0", "1", "2", "3", "4", "5", "6", "7", "8", "9",
	}
}