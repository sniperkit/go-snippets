package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func penultimate(q string) string {
	t := strings.Fields(q)
	if len(t) < 2 {
		return ""
	}
	return t[len(t)-2]
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(penultimate(scanner.Text()))
	}
}
