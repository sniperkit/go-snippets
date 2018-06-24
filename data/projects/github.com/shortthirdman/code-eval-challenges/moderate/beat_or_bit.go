package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func g2d(a uint) uint {
	a ^= a >> 4
	a ^= a >> 2
	a ^= a >> 1
	return a
}

func beatOrBit(q string) string {
	var v uint
	s := strings.Split(q, " | ")
	for i := range s {
		fmt.Sscanf(s[i], "%b", &v)
		s[i] = fmt.Sprint(g2d(v))
	}
	return strings.Join(s, " | ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(beatOrBit(scanner.Text()))
	}
}
