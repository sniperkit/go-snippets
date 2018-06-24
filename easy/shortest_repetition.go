package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func shortestRep(q string) int {
	r := 1
	for i := 1; i < len(q); i++ {
		if q[i] != q[i-r] {
			r = i + 1
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
		fmt.Println(shortestRep(scanner.Text()))
	}
}
