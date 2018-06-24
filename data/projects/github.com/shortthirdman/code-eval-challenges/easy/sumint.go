package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func sumint(n, s int) int {
	return n + s
}

func main() {
	var n, s int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &n)
		s = sumint(n, s)
	}
	fmt.Println(s)
}
