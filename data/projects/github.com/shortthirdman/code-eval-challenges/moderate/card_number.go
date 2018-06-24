package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func cardNumber(q string) bool {
	var e int
	s := strings.Split(q, "")
	t := make([]int, len(s))
	for i := range s {
		fmt.Sscan(s[i], &t[i])
	}
	for i := len(t) - 2; i >= 0; i -= 2 {
		t[i] *= 2
		if t[i] > 9 {
			t[i] = t[i]%10 + 1
		}
	}
	for i := range t {
		e += t[i]
	}
	return e%10 == 0
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
		if cardNumber(strings.Map(mapf, scanner.Text())) {
			fmt.Println("1")
		} else {
			fmt.Println("0")
		}
	}
}
