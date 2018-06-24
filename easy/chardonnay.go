package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func containsAll(a, b string) bool {
	for _, i := range b {
		ix := strings.IndexRune(a, i)
		if ix == -1 {
			return false
		}
		a = a[:ix] + a[ix+1:]
	}
	return true
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var r []string
		s := strings.Split(scanner.Text(), "|")
		t := strings.Fields(strings.TrimSpace(s[0]))
		u := strings.TrimSpace(s[1])
		for _, i := range t {
			if containsAll(i, u) {
				r = append(r, i)
			}
		}
		if len(r) == 0 {
			fmt.Println("False")
		} else {
			fmt.Println(strings.Join(r, " "))
		}
	}
}
