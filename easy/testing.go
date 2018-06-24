package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func test(s, t string) string {
	var r int
	for i := range s {
		if s[i] != t[i] {
			r++
			if r > 6 {
				break
			}
		}
	}
	switch {
	case r == 0:
		return "Done"
	case r <= 2:
		return "Low"
	case r <= 4:
		return "Medium"
	case r <= 6:
		return "High"
	default:
		return "Critical"
	}
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " | ")
		fmt.Println(test(s[0], s[1]))
	}
}
