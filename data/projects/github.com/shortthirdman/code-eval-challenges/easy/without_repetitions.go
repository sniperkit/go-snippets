package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func withoutRepetitions(p, q string) (string, string) {
	if p != q {
		return q, q
	}
	return p, ""
}

func main() {
	var p, q string
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		p, q = withoutRepetitions(p, scanner.Text())
		fmt.Print(q)
	}
	if p != "\n" {
		fmt.Println()
	}
}
