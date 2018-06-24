package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func romanArabic(q string) int {
	var n, b, r, s int
	for i := 0; i < len(q); i += 2 {
		a := int(q[i] - '0')
		switch q[i+1] {
		case 'M':
			r = 1000
		case 'D':
			r = 500
		case 'C':
			r = 100
		case 'L':
			r = 50
		case 'X':
			r = 10
		case 'V':
			r = 5
		case 'I':
			r = 1
		}
		if r > s {
			n -= b * s
		} else {
			n += b * s
		}
		b, s = a, r
	}
	return n + b*s
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(romanArabic(scanner.Text()))
	}
}
