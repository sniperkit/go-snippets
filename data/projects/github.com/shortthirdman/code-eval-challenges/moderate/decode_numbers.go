package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func decode(s string) uint {
	for len(s) > 1 {
		switch {
		case s[0] != '1' && s[0] != '2':
			s = s[1:len(s)]
		case s[1] == '0' || (s[0] == '2' && s[1] > '6'):
			s = s[2:len(s)]
		default:
			if len(s) == 2 {
				return 2
			}
			return decode(s[1:len(s)]) + decode(s[2:len(s)])
		}
	}
	return 1
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(decode(scanner.Text()))
	}
}
