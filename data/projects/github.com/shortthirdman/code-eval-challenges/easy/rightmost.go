package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func rightmost(p, q string) int {
	return strings.LastIndex(p, q)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if s := strings.Split(scanner.Text(), ","); len(s) > 1 {
			fmt.Println(rightmost(s[0], s[1]))
		}
	}
}
