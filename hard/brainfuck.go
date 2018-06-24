package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	var i int
	data, err := os.Open(os.Args[1])
	stdIn := bufio.NewReader(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var (
			p      int
			bstack []int
			f      bool
		)
		s := scanner.Text()
		b := make(map[int]int)
		for i = 0; i < len(s); i++ {
			if s[i] == '[' {
				bstack = append(bstack, i)
			} else if s[i] == ']' {
				if len(bstack) == 0 {
					f = true
					break
				}
				b[bstack[len(bstack)-1]], b[i] = i, bstack[len(bstack)-1]
				bstack = bstack[:len(bstack)-1]
			}
		}
		if f || len(bstack) > 0 {
			fmt.Println("Error: unmatched brackets")
			continue
		}
		t := []int{0}
		for i = 0; i < len(s); i++ {
			switch s[i] {
			case '+':
				t[p]++
				if t[p] == 256 {
					t[p] = 0
				}
			case '-':
				t[p]--
				if t[p] == -1 {
					t[p] = 255
				}
			case '.':
				fmt.Printf("%c", t[p])
			case ',':
				input, _ := stdIn.ReadByte()
				t[p] = int(input)
			case '<':
				if p > 0 {
					p--
				} else {
					t = append([]int{0}, t...)
				}
			case '>':
				p++
				if p == len(t) {
					t = append(t, 0)
				}
			case '[':
				if t[p] == 0 {
					i = b[i]
				}
			case ']':
				if t[p] != 0 {
					i = b[i]
				}
			}
		}
		fmt.Println()
	}
}
