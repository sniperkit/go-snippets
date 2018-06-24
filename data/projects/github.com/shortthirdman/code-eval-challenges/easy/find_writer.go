package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func findWriter(s, t string) string {
	var f int
	u := strings.Fields(t)
	r := make([]byte, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &f)
		r[ix] = byte(s[f-1])
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
		if scanner.Text() == "" {
			continue
		}
		t := strings.Split(scanner.Text(), "|")
		fmt.Println(findWriter(t[0], t[1]))
	}
}
