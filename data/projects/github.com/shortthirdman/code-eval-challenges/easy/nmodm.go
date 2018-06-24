package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func mod(n, m int) int {
	return n - (n/m)*m
}

func main() {
	var n, m int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d", &n, &m)
		fmt.Println(mod(n, m))
	}
}
