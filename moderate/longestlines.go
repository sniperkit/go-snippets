package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func longest(p []string, q string) []string {
	for ix, i := range p {
		if len(q) > len(i) {
			p = append(p[:ix], append([]string{q}, p[ix:len(p)-1]...)...)
			break
		}
	}
	return p
}

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Scan()
	fmt.Sscan(scanner.Text(), &n)
	p := make([]string, n)
	for scanner.Scan() {
		p = longest(p, scanner.Text())
	}
	fmt.Println(strings.Join(p, "\n"))
}
