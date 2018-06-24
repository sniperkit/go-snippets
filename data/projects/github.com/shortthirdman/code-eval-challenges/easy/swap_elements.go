package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func swapElements(p, q string) string {
	var a, b int
	r, s := strings.Fields(p), strings.Fields(q)
	for _, i := range s {
		fmt.Sscanf(i, "%d-%d,", &a, &b)
		r[a], r[b] = r[b], r[a]
	}
	return strings.Join(r, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), " : ")
		fmt.Println(swapElements(s[0], s[1]))
	}
}
