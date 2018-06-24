package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func multiplyLists(q string) string {
	var u, v int
	s := strings.Split(q, " | ")
	x, y := strings.Fields(s[0]), strings.Fields(s[1])
	r := make([]string, len(x))
	for i := range x {
		fmt.Sscan(x[i], &u)
		fmt.Sscan(y[i], &v)
		r[i] = fmt.Sprint(u * v)
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
		fmt.Println(multiplyLists(scanner.Text()))
	}
}
