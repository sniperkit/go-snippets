package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func stepwise(s string) string {
	var (
		maxw string
		maxl int
	)
	t := strings.Fields(s)
	for _, i := range t {
		if len(i) > maxl {
			maxw, maxl = i, len(i)
		}
	}
	r := make([]string, maxl)
	for i := 0; i < maxl; i++ {
		r[i] = strings.Repeat("*", i) + string(maxw[i])
	}
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
		fmt.Println(stepwise(scanner.Text()))
	}
}
