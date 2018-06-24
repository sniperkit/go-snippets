package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func sameend(p, q string) bool {
	return strings.HasSuffix(p, q)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		t := strings.Split(scanner.Text(), ",")
		if len(t) < 2 {
			continue
		}
		if sameend(t[0], t[1]) {
			fmt.Println(1)
		} else {
			fmt.Println(0)
		}
	}
}
