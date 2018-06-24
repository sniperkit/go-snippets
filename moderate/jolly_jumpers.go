package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func jolly(q string) bool {
	var v, vp int
	s := strings.Fields(q)[1:]
	n := len(s)
	if n == 1 {
		return true
	}
	u := make([]bool, n-1)
	fmt.Sscan(s[0], &vp)
	for i := 1; i < n; i++ {
		fmt.Sscan(s[i], &v)
		x := abs(v - vp)
		if x >= n || x == 0 || u[x-1] {
			return false
		}
		u[x-1], vp = true, v
	}
	return true
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if jolly(scanner.Text()) {
			fmt.Println("Jolly")
		} else {
			fmt.Println("Not jolly")
		}
	}
}
