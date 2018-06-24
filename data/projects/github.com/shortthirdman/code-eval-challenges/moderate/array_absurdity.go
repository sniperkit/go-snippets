package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func absurd(n int, t []string) int {
	var v int
	m := make([]bool, n-1)
	for _, i := range t {
		fmt.Sscanf(i, "%d", &v)
		if m[v] {
			return v
		}
		m[v] = true
	}
	return -1
}

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Sscanf(s[0], "%d", &n)
		fmt.Println(absurd(n, strings.Split(s[1], ",")))
	}
}
