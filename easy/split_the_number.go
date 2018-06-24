package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func split(s, t string) int {
	var n, m, p int
	var neg bool
	if strings.Contains(t, "+") {
		p = strings.Index(t, "+")
	} else {
		p = strings.Index(t, "-")
		neg = true
	}
	fmt.Sscan(s[0:p], &n)
	fmt.Sscanf(s[p:len(s)], "%d", &m) // may have leading 0
	if neg {
		return n - m
	}
	return n + m
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		fmt.Println(split(s[0], s[1]))
	}
}
