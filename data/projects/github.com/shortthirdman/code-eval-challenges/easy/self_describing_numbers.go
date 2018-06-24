package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func self(s string) bool {
	var n, r int64
	fmt.Sscan(s, &n)
	p := make([]int, len(s))
	for m := n; m > 0; m /= 10 {
		if v := int(m % 10); v < len(p) {
			p[v]++
		} else {
			return false
		}
	}
	for _, i := range p {
		r = r*10 + int64(i)
	}
	return n == r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if self(scanner.Text()) {
			fmt.Println(1)
		} else {
			fmt.Println(0)
		}
	}
}
