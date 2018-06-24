package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func swapNumbers(q string) string {
	s := strings.Fields(q)
	for ix, i := range s {
		s[ix] = string(i[len(i)-1]) + i[1:len(i)-1] + string(i[0])
	}
	return strings.Join(s, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(swapNumbers(scanner.Text()))
	}
}
