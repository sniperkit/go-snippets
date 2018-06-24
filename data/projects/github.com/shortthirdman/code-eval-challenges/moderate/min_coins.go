package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

var b []uint

func init() {
	b = []uint{0, 1, 2, 1, 2}
}

func minCoins(n uint) uint {
	return n/5 + b[n%5]
}

func main() {
	var n uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &n)
		fmt.Println(minCoins(n))
	}
}
