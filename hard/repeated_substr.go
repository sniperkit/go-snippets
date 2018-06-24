package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func repeatedSubstring(q string) string {
	l, r := len(q), ""
	if l == 0 {
		return "NONE"
	}
	for n, c := 1, 0; n <= (l-c)/2; n++ {
		var f bool
		for i := c; i < l-n; i++ {
			subs := q[i : i+n]
			if len(strings.TrimSpace(subs)) == 0 {
				continue
			}
			if strings.Contains(q[i+n:l], subs) {
				r, c, f = subs, i, true
				break
			}
		}
		if !f {
			break
		}
	}
	if r == "" {
		return "NONE"
	} else {
		return r
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
		fmt.Println(repeatedSubstring(scanner.Text()))
	}
}
