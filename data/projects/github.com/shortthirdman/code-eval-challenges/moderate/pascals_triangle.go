package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func pascal(a uint) string {
	r := []string{"1"}
	for i := uint(1); i < a; i++ {
		r = append(r, "1")
		for b, j := uint(1), uint(1); j <= i; j++ {
			b = (b * (i + 1 - j)) / j
			r = append(r, fmt.Sprint(b))
		}
	}
	return strings.Join(r, " ")
}

func main() {
	var a uint
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscan(scanner.Text(), &a)
		fmt.Println(pascal(a))
	}
}
