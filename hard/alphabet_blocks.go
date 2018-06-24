package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func alpha(b [][]int) bool {
	if len(b) == 0 {
		return true
	} else if len(b[0]) == 0 {
		return false
	}
	for _, i := range b[0] {
		c := make([][]int, len(b)-1)
		for j := 1; j < len(b); j++ {
			c[j-1] = make([]int, len(b[j]))
			copy(c[j-1], b[j])
		}
		for jx, j := range c {
			for kx, k := range j {
				if k == i {
					c[jx][kx], c[jx] = c[jx][len(c[jx])-1], c[jx][:len(c[jx])-1]
					break
				}
			}
		}
		if alpha(c) {
			return true
		}
	}
	return false
}

func contains(a string, b rune) bool {
	for _, i := range a {
		if i == b {
			return true
		}
	}
	return false
}

func alphabet(p, q string) bool {
	t := strings.Fields(q)
	n, m := len(p), len(t)
	blocks := make([][]int, n)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if contains(t[j], rune(p[i])) {
				blocks[i] = append(blocks[i], j)
			}
		}
	}
	return alpha(blocks)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " | ")
		if alphabet(s[1], s[2]) {
			fmt.Println("True")
		} else {
			fmt.Println("False")
		}
	}
}
