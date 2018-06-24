package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func even(n int) int {
	return 1 - n&1
}

func main() {
	var m int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &m)
		fmt.Println(even(m))
	}
}
