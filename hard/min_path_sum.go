package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var n int
		fmt.Sscan(scanner.Text(), &n)
		b := make([]int, n)
		scanner.Scan()
		s := strings.Split(scanner.Text(), ",")
		fmt.Sscan(s[0], &b[0])
		for i := 1; i < n; i++ {
			var v int
			fmt.Sscan(s[i], &v)
			b[i] = b[i-1] + v
		}
		for i := 1; i < n; i++ {
			scanner.Scan()
			s = strings.Split(scanner.Text(), ",")
			var v int
			fmt.Sscan(s[0], &v)
			b[0] += v
			for j := 1; j < n; j++ {
				var v int
				fmt.Sscan(s[j], &v)
				b[j] = min(b[j], b[j-1]) + v
			}
		}
		fmt.Println(b[n-1])
	}
}
