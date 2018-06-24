package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func dec2bin(a uint) string {
	return fmt.Sprintf("%b", a)
}

func main() {
	var a uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if len(scanner.Text()) > 0 {
			fmt.Sscan(scanner.Text(), &a)
			fmt.Println(dec2bin(a))
		}
	}
}
