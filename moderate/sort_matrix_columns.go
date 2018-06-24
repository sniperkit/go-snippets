package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type column []int
type columns []column

func (slice columns) Len() int { return len(slice) }
func (slice columns) Less(i, j int) bool {
	var k int
	for k < len(slice[i])-1 && slice[i][k] == slice[j][k] {
		k++
	}
	return slice[i][k] < slice[j][k]
}
func (slice columns) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func sortMatrix(q string) string {
	var i, j int
	s := strings.Split(q, " | ")
	n := strings.Count(s[0], " ") + 1
	m := make([]column, n)
	for i = 0; i < n; i++ {
		m[i] = make(column, n)
	}
	for i = 0; i < n; i++ {
		t := strings.Fields(s[i])
		for j = 0; j < n; j++ {
			fmt.Sscan(t[j], &m[j][i])
		}
	}
	sort.Sort(columns(m))
	r, t := make([]string, n), make([]string, n)
	for i = 0; i < n; i++ {
		for j = 0; j < n; j++ {
			t[j] = fmt.Sprint(m[j][i])
		}
		r[i] = strings.Join(t, " ")
	}
	return strings.Join(r, " | ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(sortMatrix(scanner.Text()))
	}
}
