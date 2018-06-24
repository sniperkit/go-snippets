package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func largestSum(q string) int {
	s := strings.Split(q, ",")
	var l, m, v int
	fmt.Sscan(strings.TrimSpace(s[0]), &m)
	for _, i := range s {
		fmt.Sscan(strings.TrimSpace(i), &v)
		if v+l > m {
			m = v + l
		}
		if v+l > 0 {
			l += v
		} else {
			l = 0
		}
	}
	return m
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(largestSum(scanner.Text()))
	}
}
