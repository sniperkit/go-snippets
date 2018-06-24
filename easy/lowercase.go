package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func toLower(q string) string {
	return strings.ToLower(q)
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(toLower(scanner.Text()))
	}
}
