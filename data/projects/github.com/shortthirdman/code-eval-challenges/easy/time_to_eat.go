package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func timeToEat(q string) string {
	s := strings.Fields(q)
	sort.Sort(sort.Reverse(sort.StringSlice(s)))
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
		fmt.Println(timeToEat(scanner.Text()))
	}
}
