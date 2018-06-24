package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func details(s string) int {
	t := strings.Split(s, ",")
	m := len(t[0])
	for _, i := range t {
		if n := strings.Index(i, "Y") - strings.LastIndex(i, "X"); n < m {
			m = n
		}
	}
	return m - 1
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(details(scanner.Text()))
	}
}
