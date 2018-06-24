package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var v map[byte]int

func init() {
	v = map[byte]int{'2': 2, '3': 3, '4': 4, '5': 5, '6': 6, '7': 7,
		'8': 8, '9': 9, '1': 10, 'J': 11, 'Q': 12, 'K': 13, 'A': 14}
}

func rank(c string, t byte) int {
	r := v[c[0]]
	if c[len(c)-1] == t {
		r += 13
	}
	return r
}

func win(l, r string, t byte) string {
	lr, rr := rank(l, t), rank(r, t)
	if lr > rr {
		return l
	} else if lr < rr {
		return r
	}
	return l + " " + r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		fmt.Println(win(s[0], s[1], s[3][0]))
	}
}
