package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func rollerCoaster(q string) string {
	var u bool
	r := make([]byte, len(q))
	for ix, i := range q {
		if (i >= 'a' && i <= 'z') || (i >= 'A' && i <= 'Z') {
			u = !u
			if u {
				r[ix] = byte(i) & 223
			} else {
				r[ix] = byte(i) | 32
			}
		} else {
			r[ix] = byte(i)
		}
	}
	return string(r)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(rollerCoaster(scanner.Text()))
	}
}
