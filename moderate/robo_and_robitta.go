package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func robo(a, b, x, y uint) uint {
	var r uint
	for b != y {
		r, a, b, x, y = r+a, b-1, a, b-y, x
	}
	return r + x
}

func main() {
	var a, b, x, y uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%dx%d | %d %d", &a, &b, &x, &y)
		fmt.Println(robo(a, b, x, y))
	}
}
