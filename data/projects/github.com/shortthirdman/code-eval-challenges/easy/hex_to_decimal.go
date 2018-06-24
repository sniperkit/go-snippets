package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func hex2dec(q string) (r int) {
	fmt.Sscanf(q, "%x", &r)
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
		fmt.Println(hex2dec(scanner.Text()))
	}
}
