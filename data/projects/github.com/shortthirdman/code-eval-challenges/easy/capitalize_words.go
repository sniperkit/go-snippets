package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func capitalize(q string) string {
	return strings.Title(q)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(capitalize(scanner.Text()))
	}
}
