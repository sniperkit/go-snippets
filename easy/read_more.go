package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func readMore(q string) string {
	if len(q) > 55 {
		for i := 39; i >= 0; i-- {
			if q[i] == ' ' {
				q = q[:i]
				break
			}
		}
		if len(q) > 55 {
			q = q[:40]
		}
		q += "... <Read More>"
	}
	return q
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(readMore(scanner.Text()))
	}
}
