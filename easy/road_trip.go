package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func roadTrip(s string) string {
	var v, c int
	t := strings.Split(s, ";")
	l, m := make([]int, len(t)-1), make([]string, len(t)-1)
	for i := range l {
		u := strings.Split(t[i], ",")
		fmt.Sscan(u[1], &v)
		l[i] = v
	}
	sort.Ints(l)
	for i := range l {
		m[i], c = fmt.Sprint(l[i]-c), l[i]
	}
	return strings.Join(m, ",")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(roadTrip(scanner.Text()))
	}
}
