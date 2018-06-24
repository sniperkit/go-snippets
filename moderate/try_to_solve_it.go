package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func deco(r rune) rune {
	switch {
	case r >= 'a' && r <= 'f':
		return r + 20
	case r >= 'g' && r <= 'm':
		return r + 7
	case r >= 'n' && r <= 't':
		return r - 7
	case r >= 'u' && r <= 'z':
		return r - 20
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
		fmt.Println(strings.Map(deco, scanner.Text()))
	}
}
