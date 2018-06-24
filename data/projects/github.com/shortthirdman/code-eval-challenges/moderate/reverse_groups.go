package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func reverseGroups(q string) string {
	var (
		m int
		r []string
	)
	t := strings.Split(q, ";")
	fmt.Sscan(t[1], &m)
	u := strings.Split(t[0], ",")
	mi := m
	for ix := range u {
		var o int
		if mi <= len(u)-ix {
			mi--
			o = mi - ix%m
			if mi == 0 {
				mi = m
			}
		}
		r = append(r, u[ix+o])
	}
	return strings.Join(r, ",")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(reverseGroups(scanner.Text()))
	}
}
