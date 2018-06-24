package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func upper(b byte) byte {
	return b & 223
}

func mask(s, t string) string {
	r := make([]byte, len(s))
	for ix, i := range s {
		if t[ix] == '1' {
			r[ix] = upper(byte(i))
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
		f := len(scanner.Text()) / 2
		fmt.Println(mask(scanner.Text()[:f], scanner.Text()[f+1:]))
	}
}
