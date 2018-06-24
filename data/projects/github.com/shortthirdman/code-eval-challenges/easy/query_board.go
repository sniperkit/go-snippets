package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func main() {
	var p, v int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	m := make([][]int, 256)
	for i := 0; i < 256; i++ {
		m[i] = make([]int, 256)
	}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		t := strings.Fields(scanner.Text())
		if strings.HasPrefix(t[0], "Query") {
			var sum int
			fmt.Sscan(t[1], &p)
			if strings.HasSuffix(t[0], "Row") {
				for i := 0; i < 256; i++ {
					sum += m[p][i]
				}
			} else {
				for i := 0; i < 256; i++ {
					sum += m[i][p]
				}
			}
			fmt.Println(sum)
		} else {
			fmt.Sscan(t[1], &p)
			fmt.Sscan(t[2], &v)
			if strings.HasSuffix(t[0], "Row") {
				for i := 0; i < 256; i++ {
					m[p][i] = v
				}
			} else {
				for i := 0; i < 256; i++ {
					m[i][p] = v
				}
			}
		}
	}
}
