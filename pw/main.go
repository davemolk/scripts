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
	rand.Seed(time.Now().UnixNano())
	ch, err := read()
	if err != nil {
		log.Fatal(err)
	}

	pw := checkBoxes(specials(nums(lower(upper(ch)))))
	fmt.Println(pw)	
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

func upper(str <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for char := range str {
			if rand.Intn(2) == 1 {
				char = strings.ToUpper(char)
			}
			ch <- char
		}
	}()
	return ch
}

func lower(str <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		for char := range str {
			if rand.Intn(2) == 1 {
				char = strings.ToLower(char)
			}
			ch <- char
		}
	}()
	return ch
}

func nums(str <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		nums := "0123456789"
		for char := range str {
			if rand.Intn(2) == 1 {
				char = string(nums[rand.Intn(len(nums))])
			}
			ch <- char
		}
	}()
	return ch
}

func specials(str <-chan string) <-chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		sp := `" !#$%&()*,-./:;<=>?@[\]^_{|}~`
		for char := range str {
			if rand.Intn(2) == 1 {
				char = string(sp[rand.Intn(len(sp))])
			}
			ch <- char
		}
	}()
	return ch
}

func checkBoxes(str <-chan string) string {
	var preCheck []string
	for char := range str {
		preCheck = append(preCheck, char)
	}
	nums := "0123456789"
	specials := `" !#$%&()*,-./:;<=>?@[]^_{|}~`
	lowers := "abcdefghijlkmnopqrstuvwxyz"
	uppers := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	buf := make([]byte, 4)
	buf[0] = nums[rand.Intn(len(nums))]
	buf[1] = specials[rand.Intn(len(specials))]
	buf[2] = lowers[rand.Intn(len(lowers))]
	buf[3] = uppers[rand.Intn(len(uppers))]

	for _, b := range buf {
		preCheck = append(preCheck, string(b))
	}
	rand.Shuffle(len(preCheck), func(i, j int) {
		preCheck[i], preCheck[j] = preCheck[j], preCheck[i]
	})

	pw := strings.Replace(strings.Join(preCheck, ""), " ", "", -1)
	
	return pw
}