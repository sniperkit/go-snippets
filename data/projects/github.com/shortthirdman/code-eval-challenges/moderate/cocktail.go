package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func cocktail(q string, n int) string {
	u := strings.Fields(q)
	n = min(n, len(u)/2)
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &v[ix])
	}
	for i := 0; i < n; i++ {
		for j := i + 1; j < len(v)-i; j++ {
			if v[j-1] > v[j] {
				v[j-1], v[j] = v[j], v[j-1]
			}
		}
		for j := len(v) - i - 2; j >= i+1; j-- {
			if v[j-1] > v[j] {
				v[j-1], v[j] = v[j], v[j-1]
			}
		}
	}
	for ix, i := range v {
		u[ix] = fmt.Sprint(i)
	}
	return strings.Join(u, " ")
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
		t := strings.Split(scanner.Text(), " | ")
		fmt.Sscan(t[1], &n)
		fmt.Println(cocktail(t[0], n))
	}
}
