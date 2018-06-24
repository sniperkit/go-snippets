package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func max(a, b int) int {
	if b > a {
		return b
	}
	return a
}

func highestScore(s string) string {
	var (
		a int
		r []int
		u []string
	)
	for ix, i := range strings.Split(s, "|") {
		for jx, j := range strings.Fields(i) {
			fmt.Sscan(j, &a)
			if ix == 0 {
				r = append(r, a)
			} else {
				r[jx] = max(r[jx], a)
			}
		}
	}
	for _, i := range r {
		u = append(u, fmt.Sprint(i))
	}
	return strings.Join(u, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(highestScore(scanner.Text()))
	}
}
