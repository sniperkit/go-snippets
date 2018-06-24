package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var moy map[string]int

func init() {
	moy = map[string]int{"Jan": 0, "Feb": 1, "Mar": 2, "Apr": 3,
		"May": 4, "Jun": 5, "Jul": 6, "Aug": 7,
		"Sep": 8, "Oct": 9, "Nov": 10, "Dec": 11}
}

func month(s string) int {
	var k int
	fmt.Sscan(s[4:8], &k)
	return 12*(k-1990) + moy[s[0:3]]
}

func working(q string) int {
	work := make(map[int]bool)
	t := strings.Split(q, "; ")
	for _, i := range t {
		t0, t1 := month(i[0:8]), month(i[9:17])
		for j := t0; j <= t1; j++ {
			work[j] = true
		}
	}
	return len(work) / 12
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(working(scanner.Text()))
	}
}
