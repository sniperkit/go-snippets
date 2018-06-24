package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
	"strings"
)

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func diagonal(c, h []int) bool {
	a, b, p, q := float64(h[0]), float64(h[1]), float64(c[0]), float64(c[1])
	return (p > a) && (b >= (2*p*q*a+(p*p-q*q)*math.Sqrt(p*p+q*q-a*a))/(p*p+q*q))
}

func pile(a, b, c, d int, q string) string {
	var e, f, j int
	g, h := make([]int, 2), make([]int, 3)
	g[0], g[1] = abs(a-c), abs(b-d)
	sort.Ints(g)
	var r []int
	for _, i := range strings.Split(q, ";") {
		fmt.Sscanf(i, "(%d [%d,%d,%d] [%d,%d,%d])", &j, &a, &b, &c, &d, &e, &f)
		h[0], h[1], h[2] = abs(a-d), abs(b-e), abs(c-f)
		sort.Ints(h)
		if h[0] <= g[0] && h[1] <= g[1] { // || diagonal(h, g)
			r = append(r, j)
		}
	}
	if len(r) == 0 {
		return "-"
	}
	sort.Ints(r)
	var r2 []string
	for _, i := range r {
		r2 = append(r2, fmt.Sprint(i))
	}
	return strings.Join(r2, ",")
}

func main() {
	var a, b, c, d int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		t := strings.Split(scanner.Text(), "|")
		fmt.Sscanf(t[0], "[%d,%d] [%d,%d]", &a, &b, &c, &d)
		fmt.Println(pile(a, b, c, d, t[1]))
	}
}
