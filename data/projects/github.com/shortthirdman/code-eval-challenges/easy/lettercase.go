package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func lcase(s string) (float64, float64) {
	var c, l float64
	for _, i := range s {
		if i >= 'a' && i <= 'z' {
			l++
			c++
		} else if i >= 'A' && i <= 'Z' {
			c++
		}
	}
	if c == 0 {
		return 0, 0
	}
	r := 100 * l / c
	return r, 100 - r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		l, u := lcase(scanner.Text())
		fmt.Printf("lowercase: %.2f uppercase: %.2f\n", l, u)
	}
}
