package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func longest(q string) (r string) {
	for _, i := range strings.Fields(q) {
		if len(i) > len(r) {
			r = i
		}
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(longest(scanner.Text()))
	}
}
