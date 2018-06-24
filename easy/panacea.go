package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func panacea(p, q string) bool {
	var v, m int
	l, r := strings.Fields(p), strings.Fields(q)
	for _, i := range l {
		fmt.Sscanf(i, "%x", &v)
		m += v
	}
	for _, i := range r {
		fmt.Sscanf(i, "%b", &v)
		m -= v
	}
	return m <= 0
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " | ")
		if panacea(s[0], s[1]) {
			fmt.Println("True")
		} else {
			fmt.Println("False")
		}
	}
}
