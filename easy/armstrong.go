package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func powi(y, x int) int {
	ret := 1
	for x > 0 {
		if x&1 > 0 {
			ret *= y
		}
		y *= y
		x >>= 1
	}
	return ret
}

func armstrong(n int) bool {
	var e, t int
	for a := n; a > 0; a /= 10 {
		e++
	}
	for a := n; a > 0; a /= 10 {
		t += powi(a%10, e)
	}
	return n == t
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
		fmt.Sscan(scanner.Text(), &n)
		if armstrong(n) {
			fmt.Println("True")
		} else {
			fmt.Println("False")
		}
	}
}
