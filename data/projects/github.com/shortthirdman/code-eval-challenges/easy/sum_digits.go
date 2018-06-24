package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func sumDigits(s string) (r uint) {
	for _, i := range s {
		if i >= '0' && i <= '9' {
			r += uint(i - '0')
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
		fmt.Println(sumDigits(scanner.Text()))
	}
}
