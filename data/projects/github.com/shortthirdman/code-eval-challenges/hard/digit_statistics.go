package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

var stats [][]int

func init() {
	stats = [][]int{[]int{2, 4, 8, 6}, []int{3, 9, 7, 1},
		[]int{4, 6}, []int{5}, []int{6},
		[]int{7, 9, 3, 1}, []int{8, 4, 2, 6},
		[]int{9, 1}}
}

func digitStatistics(a, n int) string {
	a -= 2
	m := n / len(stats[a])
	q, r := make([]int, 10), make([]string, 10)
	for _, i := range stats[a] {
		q[i] += m
	}
	for i := 0; i < n%len(stats[a]); i++ {
		q[stats[a][i]]++
	}
	for i := 0; i < 10; i++ {
		r[i] = fmt.Sprintf("%d: %d", i, q[i])
	}
	return strings.Join(r, ", ")
}

func main() {
	var a, n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d %d", &a, &n)
		fmt.Println(digitStatistics(a, n))
	}
}
