package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func pos(l, r uint8) string {
	return string([]byte{'a' + l, '1' + r})
}

func knight(q string) string {
	var m []string
	l, r := q[0]-'a', q[1]-'1'
	if l > 1 {
		if r > 0 {
			m = append(m, pos(l-2, r-1))
		}
		if r < 7 {
			m = append(m, pos(l-2, r+1))
		}
	}
	if l > 0 {
		if r > 1 {
			m = append(m, pos(l-1, r-2))
		}
		if r < 6 {
			m = append(m, pos(l-1, r+2))
		}
	}
	if l < 7 {
		if r > 1 {
			m = append(m, pos(l+1, r-2))
		}
		if r < 6 {
			m = append(m, pos(l+1, r+2))
		}
	}
	if l < 6 {
		if r > 0 {
			m = append(m, pos(l+2, r-1))
		}
		if r < 7 {
			m = append(m, pos(l+2, r+1))
		}
	}
	return strings.Join(m, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(knight(scanner.Text()))
	}
}
