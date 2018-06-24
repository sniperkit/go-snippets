package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func hiddenDigits(q string) string {
	maph := func(r rune) rune {
		if strings.Contains("0123456789", string(r)) {
			return r
		} else if strings.Contains("abcdefghij", string(r)) {
			return r - 'a' + '0'
		}
		return -1
	}
	if r := strings.Map(maph, q); len(r) > 0 {
		return r
	} else {
		return "NONE"
	}
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(hiddenDigits(scanner.Text()))
	}
}
