package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func nextPermu(s string) (string, bool) {
	k, l := -1, 0
	for i := len(s) - 2; i >= 0; i-- {
		if s[i] < s[i+1] {
			k = i
			break
		}
	}
	if k == -1 {
		return s, true
	}
	for i := len(s) - 1; i > 0; i-- {
		if s[i] > s[k] {
			l = i
			break
		}
	}
	s = strings.Join([]string{s[0:k], string(s[l]), s[k+1 : l], string(s[k]), s[l+1 : len(s)]}, "")
	l = len(s) - k - 1
	for i := 0; l-1-i > i; i++ {
		s = strings.Join([]string{s[0 : k+1+i], string(s[len(s)-1-i]), s[k+2+i : len(s)-1-i], string(s[k+1+i]), s[len(s)-i : len(s)]}, "")
	}
	return s, false
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), "")
		sort.Strings(s)
		c := strings.Join(s, "")
		fmt.Print(c)
		var done bool
		for c, done = nextPermu(c); !done; c, done = nextPermu(c) {
			fmt.Printf(",%s", c)
		}
		fmt.Println()
	}
}
