package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

func intersect(s, t string) string {
	var (
		a    int
		x, y []int
		r    []string
	)
	u, v := strings.Split(s, ","), strings.Split(t, ",")
	for _, i := range u {
		fmt.Sscan(i, &a)
		x = append(x, a)
	}
	for _, i := range v {
		fmt.Sscan(i, &a)
		y = append(y, a)
	}
	for i, j := 0, 0; i < len(u) && j < len(v); {
		if x[i] == y[j] {
			r = append(r, fmt.Sprint(x[i]))
			i++
			j++
		} else if x[i] > y[j] {
			j++
		} else {
			i++
		}
	}
	return strings.Join(r, ",")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ";")
		fmt.Println(intersect(s[0], s[1]))
	}
}
