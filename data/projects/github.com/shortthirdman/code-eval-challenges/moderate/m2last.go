package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func m2last(q string) string {
	t := strings.Fields(q)
	var m uint
	n := uint(len(t))
	fmt.Sscan(t[n-1], &m)
	if m < n {
		return t[n-1-m]
	}
	return ""
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if r := m2last(scanner.Text()); len(r) > 0 {
			fmt.Println(r)
		}
	}
}
