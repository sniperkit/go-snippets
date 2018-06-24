package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func major(s string) string {
	var k, maxnum, maxcount int
	num := make([]int, 101)
	t := strings.Split(s, ",")
	for ix, i := range t {
		fmt.Sscan(i, &k)
		num[k]++
		if num[k] > maxcount {
			maxnum, maxcount = k, num[k]
		}
		if len(t)/2+1 > maxcount+len(t)-ix-1 {
			break
		} else if maxcount >= len(t)/2+1 {
			return fmt.Sprint(maxnum)
		}
	}
	return "None"
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(major(scanner.Text()))
	}
}
