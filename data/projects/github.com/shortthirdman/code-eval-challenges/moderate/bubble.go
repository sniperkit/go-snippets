package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func bubble(q string) string {
	var n, k int
	t := strings.Split(q, " | ")
	fmt.Sscan(t[1], &n)
	u := strings.Fields(t[0])
	v := make([]int, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &k)
		v[ix] = k
	}
	for i := 0; i < len(v)-n; i++ {
		sort.Ints(v[i : i+n+1])
	}
	for ix, i := range v {
		u[ix] = fmt.Sprint(i)
	}
	return strings.Join(u, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(bubble(scanner.Text()))
	}
}
