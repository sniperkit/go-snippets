package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func bitEqual(n, p, q uint) bool {
	return ((n & (1 << (p - 1))) == 0) == ((n & (1 << (q - 1))) == 0)
}

func main() {
	var n, p, q uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d,%d,%d", &n, &p, &q)
		fmt.Println(bitEqual(n, p, q))
	}
}
