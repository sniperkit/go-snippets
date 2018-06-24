package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func nrc(q string) string {
	s := make(map[rune]int)
	for _, i := range q {
		s[i]++
	}
	for _, i := range q {
		if s[i] == 1 {
			return string(i)
		}
	}
	return ""
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(nrc(scanner.Text()))
	}
}
