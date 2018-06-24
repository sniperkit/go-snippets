package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func blackCard(s []string, m int) string {
	for len(s) > 1 {
		n := m%len(s) - 1
		if n == -1 {
			s = s[:len(s)-1]
		} else {
			s = append(s[:n], s[n+1:]...)
		}
	}
	return s[0]
}

func main() {
	var m int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		t := strings.Split(scanner.Text(), " | ")
		fmt.Sscan(t[1], &m)
		fmt.Println(blackCard(strings.Fields(t[0]), m))
	}
}
