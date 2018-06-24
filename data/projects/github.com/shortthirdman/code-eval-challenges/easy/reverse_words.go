package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func reverse(s string) string {
	var (
		c []byte
		t []string
	)
	for _, i := range s {
		if i == ' ' {
			t, c = append([]string{string(c)}, t...), []byte{}
		} else {
			c = append(c, byte(i))
		}
	}
	t = append([]string{string(c)}, t...)
	return strings.Join(t, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if len(scanner.Text()) == 0 {
			continue
		}
		fmt.Println(reverse(scanner.Text()))
	}
}
