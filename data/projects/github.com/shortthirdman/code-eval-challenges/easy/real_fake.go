package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func realFake(q string) bool {
	s := strings.Split(q, "")
	t := make([]int, len(s))
	for i := range s {
		fmt.Sscan(s[i], &t[i])
	}
	for i := len(t) - 2; i >= 0; i -= 2 {
		t[i] *= 2
	}
	var su int
	for i := range t {
		su += t[i]
	}
	return su%10 == 0
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	mapf := func(r rune) rune {
		if r >= '0' && r <= '9' {
			return r
		}
		return -1
	}
	for scanner.Scan() {
		if realFake(strings.Map(mapf, scanner.Text())) {
			fmt.Println("Real")
		} else {
			fmt.Println("Fake")
		}
	}
}
