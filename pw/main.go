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
		nums := getNums()
		for char := range str {
			if rand.Intn(2) == 1 {
				char = nums[rand.Intn(len(nums))]
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
		sp := getSpecials()
		for char := range str {
			if rand.Intn(2) == 1 {
				char = sp[rand.Intn(len(sp))]
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
	lowers := getLetters()
	newLower := lowers[rand.Intn(len(lowers))]
	newUpper := strings.ToUpper(lowers[rand.Intn(len(lowers))])
	specials := getSpecials()
	newSpecial := specials[rand.Intn(len(specials))]
	nums := getNums()
	newNum := nums[rand.Intn(len(nums))]
	preCheck = append(preCheck, newLower, newUpper, newSpecial, newNum)

	rand.Shuffle(len(preCheck), func(i, j int) {
		preCheck[i], preCheck[j] = preCheck[j], preCheck[i]
	})

	noSpaces := strings.Replace(strings.Join(preCheck, ""), " ", "", -1)
	// cause that one time...
	noReturn := strings.Replace(noSpaces, "\n", `n\`, -1)
	
	return noReturn
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

func getLetters() []string{
	return []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
}