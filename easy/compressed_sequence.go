package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func compress(s string) string {
	var (
		r []string
		n string
		c int
	)
	for _, i := range strings.Fields(s) {
		if i == n {
			c++
			continue
		}
		if c > 0 {
			r = append(append(r, fmt.Sprint(c)), n)
		}
		c, n = 1, i
	}
	r = append(append(r, fmt.Sprint(c)), n)
	return strings.Join(r, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(compress(scanner.Text()))
	}
}
