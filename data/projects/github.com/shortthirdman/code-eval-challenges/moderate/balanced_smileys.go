package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func isBalanced(s string, c int) bool {
	for ; s != "" && c >= 0; s = s[1:] {
		switch s[0] {
		case '(':
			c++
		case ')':
			c--
		case ':':
			if s[1] == '(' {
				return isBalanced(s[2:], c) || isBalanced(s[2:], c+1)
			} else if s[1] == ')' {
				return isBalanced(s[2:], c) || isBalanced(s[2:], c-1)
			}
		}
	}
	return c == 0
}

func trimRight(s string) string {
	r := len(s)
	for ; r > 0 && s[r-1] != '(' && s[r-1] != ')'; r-- {
	}
	return s[:r]
}

func checkChars(s string) bool {
	for _, i := range s {
		if i != ' ' && i != '(' && i != ')' && i != ':' && (i < 'a' || i > 'z') {
			return false
		}
	}
	return true
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if s := scanner.Text(); checkChars(s) && isBalanced(trimRight(s), 0) {
			fmt.Println("YES")
		} else {
			fmt.Println("NO")
		}
	}
}
