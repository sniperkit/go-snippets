package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func isAlpha(a rune) bool {
	return (a >= 'a' && a <= 'z') || (a >= 'A' && a <= 'Z')
}

func clean(q string) string {
	var (
		c []byte
		r []string
	)
	w := false
	for _, i := range q {
		if isAlpha(i) {
			w = true
			c = append(c, byte(i)|32)
		} else if w {
			w = false
			if len(c) > 0 {
				r = append(r, string(c))
				c = []byte{}
			}
		}
	}
	if len(c) > 0 {
		r = append(r, string(c))
	}
	return strings.Join(r, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(clean(scanner.Text()))
	}
}
