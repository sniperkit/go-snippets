package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func locks(n, m int) int {
	if n == 0 || m == 0 {
		return n
	} else if m == 1 {
		return n - 1
	}
	l := make([]bool, n)
	for i := 0; i <= m%2; i++ {
		for j := 1; j < n; j += 2 {
			l[j] = true
		}
		for j := 2; j < n; j += 3 {
			l[j] = !l[j]
		}
	}
	l[n-1] = !l[n-1]
	var c int
	for i := 0; i < n; i++ {
		if !l[i] {
			c++
		}
	}
	return c
}

func main() {
	var n, m int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d %d", &n, &m)
		fmt.Println(locks(n, m))
	}
}
