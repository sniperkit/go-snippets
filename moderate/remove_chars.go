package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func removeChars(p, q string) string {
	mapf := func(r rune) rune {
		if strings.Contains(q, string(r)) {
			return -1
		}
		return r
	}
	return strings.Map(mapf, p)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ", ")
		fmt.Println(removeChars(s[0], s[1]))
	}
}
