package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func minDist(q string) (r int) {
	var n int
	s := strings.Fields(q)
	fmt.Sscan(s[0], &n)
	t := make([]int, n)
	for i := 1; i <= n; i++ {
		fmt.Sscan(s[i], &t[i-1])
	}
	sort.Ints(t)
	for i, u := 0, t[n/2]; i < n; i++ {
		r += abs(t[i] - u)
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(minDist(scanner.Text()))
	}
}
