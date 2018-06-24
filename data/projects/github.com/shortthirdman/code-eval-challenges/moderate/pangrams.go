package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func pangram(q string) string {
	var r []byte
	q = strings.ToLower(q)
	for i := byte('a'); i <= 'z'; i++ {
		if !strings.Contains(q, string(i)) {
			r = append(r, i)
		}
	}
	if len(r) == 0 {
		return "NULL"
	}
	return string(r)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(pangram(scanner.Text()))
	}
}
