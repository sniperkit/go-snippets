package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func kernighan(a int) (r int) {
	for ; a > 0; a &= a - 1 {
		r++
	}
	return r
}

func main() {
	var a int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &a)
		fmt.Println(kernighan(a))
	}
}
