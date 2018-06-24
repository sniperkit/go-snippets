package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func pol(o byte, a, b float32) float32 {
	switch o {
	case '*':
		return a * b
	case '/':
		return a / b
	default: // '+'
		return a + b
	}
}

func prefix(q string) int {
	var r, v float32
	s := strings.Fields(q)
	n := len(s)
	fmt.Sscan(s[n/2], &r)
	for i := 1; i <= n/2; i++ {
		fmt.Sscan(s[n/2+i], &v)
		r = pol(s[n/2-i][0], r, v)
	}
	return int(r + 0.001)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(prefix(scanner.Text()))
	}
}
