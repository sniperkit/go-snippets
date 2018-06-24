package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var digits map[string]string

func init() {
	digits = map[string]string{"zero": "0", "one": "1", "two": "2",
		"three": "3", "four": "4", "five": "5", "six": "6",
		"seven": "7", "eight": "8", "nine": "9"}
}

func wordToDigit(q string) string {
	t := strings.Split(q, ";")
	r := make([]string, len(t))
	for ix, i := range t {
		r[ix] = digits[i]
	}
	return strings.Join(r, "")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(wordToDigit(scanner.Text()))
	}
}
