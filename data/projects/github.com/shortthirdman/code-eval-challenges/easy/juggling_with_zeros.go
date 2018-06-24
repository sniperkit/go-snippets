package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func juggle(q string) (r uint64) {
	t := strings.Fields(q)
	for i := 0; i < len(t); i += 2 {
		r <<= uint(len(t[i+1]))
		if t[i] == "00" {
			r += (1 << uint(len(t[i+1]))) - 1
		}
	}
	return r
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(juggle(scanner.Text()))
	}
}
