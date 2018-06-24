package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func burrowsWheeler(s string) string {
	e, r := make([]int, len(s)), make([]byte, len(s))
	var k int
	for i := 0; i < len(s); i++ {
		e[i] = int(s[i])<<16 + i
		if s[i] == '$' {
			k = i
		}
	}
	sort.Ints(e)
	for i := 0; i < len(s); i, k = i+1, e[k]%(1<<16) {
		r[i] = byte(e[k] >> 16)
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
		fmt.Println(burrowsWheeler(strings.TrimRight(scanner.Text(), "|")))
	}
}
