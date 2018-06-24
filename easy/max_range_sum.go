package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func maxRangeSum(n int, q string) (r int) {
	t := strings.Fields(q)
	var c int
	u := make([]int, len(t))
	for ix, i := range t {
		fmt.Sscan(i, &u[ix])
	}
	for i := 0; i < n; i++ {
		c += u[i]
	}
	if c > r {
		r = c
	}
	for len(u) > n {
		c = c - u[0] + u[n]
		if c > r {
			r = c
		}
		u = u[1:]
	}
	return r
}

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Sscan(s[0], &n)
		fmt.Println(maxRangeSum(n, s[1]))
	}
}
