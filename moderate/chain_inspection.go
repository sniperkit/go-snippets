package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func chain(q string) bool {
	var (
		been []string
		f    bool
	)
	s := strings.Split(q, ";")
	m, b := make(map[string]string, len(s)), "BEGIN"
	for _, i := range s {
		t := strings.Split(i, "-")
		if t[0] == t[1] {
			return false
		}
		m[t[0]] = t[1]
	}
	for _ = range s {
		if b, f = m[b]; !f {
			return false
		}
		for _, j := range been {
			if b == j {
				return false
			}
		}
		been = append(been, b)
	}
	if b == "END" {
		return true
	}
	return false
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if chain(scanner.Text()) {
			fmt.Println("GOOD")
		} else {
			fmt.Println("BAD")
		}
	}
}
