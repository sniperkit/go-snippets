package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func subs(s, t string) uint {
	if len(t) == 0 {
		return 1
	}
	for len(s) > 0 {
		if s[0] == t[0] {
			return subs(s[1:len(s)], t[1:len(t)]) + subs(s[1:len(s)], t)
		}
		s = s[1:len(s)]
	}
	return 0
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ",")
		fmt.Println(subs(s[0], s[1]))
	}
}
