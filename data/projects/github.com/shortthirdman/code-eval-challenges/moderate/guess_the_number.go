package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func guessNumber(q string) int {
	s := strings.Fields(q)
	var c, lo, hi int
	fmt.Sscan(s[0], &hi)
	s = s[1 : len(s)-1]
	for _, i := range s {
		if i == "Lower" {
			hi = (lo+hi)/2 + c - 1
		} else {
			lo = (lo+hi)/2 + c + 1
		}
		c = (lo + hi) % 2
	}
	return (lo+hi)/2 + c
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(guessNumber(scanner.Text()))
	}
}
