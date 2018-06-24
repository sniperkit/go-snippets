package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func subcheck(s0, s1 string, q, r int) bool {
	for q < len(s0) && r < len(s1) {
		if s1[r] == '*' {
			var m bool
			for i := q; i < len(s0); i++ {
				if subcheck(s0, s1, i, r+1) {
					m = true
					break
				}
			}
			return m
		} else if s1[r] == '\\' && r < len(s1)+1 && s1[r+1] == '*' {
			if s0[q] != '*' {
				return false
			}
			q, r = q+1, r+2
		}
		if s0[q] != s1[r] {
			return false
		}
		q, r = q+1, r+1
	}
	return r >= len(s1)
}

func substring(p, q string) bool {
	q = strings.TrimLeft(q, "*") // stars without item. ignore.
	for len(q) > 1 && q[len(q)-1] == '*' && q[len(q)-2] != '\\' {
		q = q[0 : len(q)-1]
	} // redundant terminal stars.
	for i := 0; i < len(p); i++ {
		if subcheck(p, q, i, 0) {
			return true
		}
	}
	return false
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		fmt.Println(substring(s[0], s[1]))
	}
}
