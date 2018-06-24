package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var m []int

func init() {
	m = []int{1, 5, 10, 25, 50}
}

func alter(n, c int) int {
	if n < m[1] || c == 0 {
		return 1
	}
	if m[c] > n {
		return alter(n, c-1)
	}
	return alter(n-m[c], c) + alter(n, c-1)
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
		fmt.Sscanf(scanner.Text(), "%d", &n)
		fmt.Println(alter(n, 4))
	}
}
