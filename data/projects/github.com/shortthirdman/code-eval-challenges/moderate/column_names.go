package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func columnNames(c int) (r string) {
	for c > 0 {
		c--
		r = string('A'+c%26) + r
		c /= 26
	}
	return r
}

func main() {
	var c int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &c)
		fmt.Println(columnNames(c))
	}
}
