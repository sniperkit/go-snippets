package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func lwc(f []string, c string) int {
	m := len(c) / 2
	if len(f) == 0 {
		return m
	}
	for ix, i := range f {
		if len(c) == 0 || i[1] == c[0] {
			f2 := make([]string, len(f))
			copy(f2, f)
			f2[ix], f2[len(f2)-1], f2 = f2[len(f2)-1], "", f2[:len(f2)-1]
			l := lwc(f2, i+c)
			if l > m {
				m = l
			}
		}
	}
	return m
}

func wordChain(q string) string {
	s := strings.Split(q, ",")
	if len(s) < 2 {
		return "None"
	}
	for ix, i := range s {
		s[ix] = string(i[0]) + string(i[len(i)-1])
	}
	if l := lwc(s, ""); l == 1 {
		return "None"
	} else {
		return fmt.Sprint(l)
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
		fmt.Println(wordChain(scanner.Text()))
	}
}
