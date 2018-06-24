package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func justi(s string) string {
	if len(s) == 80 {
		return s
	}
	ws := strings.Fields(s)
	l := 80 - len(strings.Join(ws, ""))
	for i := 0; i < l%(len(ws)-1); i++ {
		ws[i] = ws[i] + " "
	}
	return strings.Join(ws, strings.Repeat(" ", l/(len(ws)-1)))
}

func maxwords(s string) int {
	if len(s) > 81 {
		return maxwords(s[:81])
	}
	return strings.LastIndex(s, " ")
}

func justif(s string) []string {
	t := strings.TrimLeft(s, " ")
	if len(t) <= 80 {
		return []string{t}
	}
	m := maxwords(t)
	return append([]string{justi(t[:m])}, justif(t[m:])...)
}

func justify(s string) string {
	if len(s) <= 80 {
		return s
	}
	return strings.Join(justif(s), "\n")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(justify(scanner.Text()))
	}
}
