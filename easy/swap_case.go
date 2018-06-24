package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func swapCase(c byte) byte {
	switch {
	case c >= 'a' && c <= 'z':
		return c & 223
	case c >= 'A' && c <= 'Z':
		return c | 32
	}
	return c
}

func main() {
	var c byte = '\n'
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanBytes)
	for scanner.Scan() {
		c = swapCase(scanner.Text()[0])
		fmt.Printf("%c", c)
	}
	if c != '\n' {
		fmt.Println()
	}
}
