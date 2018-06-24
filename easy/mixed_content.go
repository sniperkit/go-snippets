package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func mixedContent(q string) (r string) {
	var t, r1, r2 []string
	t = strings.Split(q, ",")
	for _, i := range t {
		_, err := strconv.ParseUint(i, 10, 64)
		if err != nil {
			r1 = append(r1, i)
		} else {
			r2 = append(r2, i)
		}
	}
	if len(r1) > 0 {
		r = strings.Join(r1, ",")
		if len(r2) > 0 {
			r += "|"
		}
	}
	if len(r2) > 0 {
		r += strings.Join(r2, ",")
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
		fmt.Println(mixedContent(scanner.Text()))
	}
}
