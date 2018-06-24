package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func lowestUnique(s string) uint {
	var k uint
	m := make([]int, 9)
	for ix, i := range strings.Fields(s) {
		fmt.Sscan(i, &k)
		if m[k-1] == 0 {
			m[k-1] = ix + 1
		} else {
			m[k-1] = -1
		}
	}
	for i := uint(0); i < 9; i++ {
		if m[i] > 0 {
			return i
		}
	}
	return 0
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(lowestUnique(scanner.Text()))
	}
}
