package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func dupes(q string) string {
	var (
		p string
		r []string
	)
	t := strings.Split(q, ",")
	for _, i := range t {
		if i != p {
			r = append(r, i)
			p = i
		}
	}
	return strings.Join(r, ",")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(dupes(scanner.Text()))
	}
}
