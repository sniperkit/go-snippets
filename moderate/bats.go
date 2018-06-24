package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func bats(l, d, n int, s []string) (c int) {
	var t int
	tx := 6 - d
	for i := 6; i <= l-6; i += d {
		if i > tx-d {
			i = tx
			if t == n {
				tx = l - 6 + d
			} else {
				fmt.Sscanf(s[t], "%d", &tx)
				t++
			}
		} else {
			c++
		}
	}
	return c
}

func main() {
	var l, d, n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d %d %d", &l, &d, &n)
		s := strings.Fields(scanner.Text())[3:]
		fmt.Println(bats(l, d, n, s))
	}
}
