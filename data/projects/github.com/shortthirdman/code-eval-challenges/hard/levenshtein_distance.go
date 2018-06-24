package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func contain(a []string, b string) bool {
	for _, i := range a {
		if i == b {
			return true
		}
	}
	return false
}

func leven(a, b string) bool {
	var d bool
	for ix, jx := 0, 0; ix < len(a); ix++ {
		if jx >= len(b) || b[jx] != a[ix] {
			if d {
				return false
			}
			d = true
			if len(a) == len(b) {
				jx++
			}
		} else {
			jx++
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
	var test []string
	scanner := bufio.NewScanner(data)
	for scanner.Scan(); scanner.Text() != "END OF INPUT"; scanner.Scan() {
		test = append(test, scanner.Text())
	}
	words := make(map[int][]string)
	for scanner.Scan() {
		words[len(scanner.Text())] = append(words[len(scanner.Text())], scanner.Text())
	}
	for _, i := range test {
		var (
			been []string
			n    int
			c    string
		)
		curr := []string{i}
		for len(curr) > 0 {
			c, curr, been, n = curr[0], curr[1:], append(been, curr[0]), n+1
			if len(c) > 1 {
				for _, j := range words[len(c)-1] {
					if !contain(been, j) && !contain(curr, j) && leven(c, j) {
						curr = append(curr, j)
					}
				}
			}
			for _, j := range words[len(c)] {
				if !contain(been, j) && !contain(curr, j) && leven(c, j) {
					curr = append(curr, j)
				}
			}
			for _, j := range words[len(c)+1] {
				if !contain(been, j) && !contain(curr, j) && leven(j, c) {
					curr = append(curr, j)
				}
			}
		}
		fmt.Println(n)
	}
}
