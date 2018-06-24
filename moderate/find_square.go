package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type point struct {
	x, y int
}

func dist2(a, b point) int {
	x, y := a.x-b.x, a.y-b.y
	return x*x + y*y
}

func findSquare(q string) bool {
	var (
		x, y int
		d    []int
	)
	t := strings.Split(q, ", ")
	u := make([]point, 4)
	for ix, i := range t {
		fmt.Sscanf(i, "(%d,%d)", &x, &y)
		u[ix] = point{x, y}
	}
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 4; j++ {
			d = append(d, dist2(u[i], u[j]))
		}
	}
	sort.Ints(d)
	return d[0] == d[1] && d[0] == d[2] && d[0] == d[3] &&
		d[4] == d[5] && 2*d[0] == d[4]
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(findSquare(scanner.Text()))
	}
}
