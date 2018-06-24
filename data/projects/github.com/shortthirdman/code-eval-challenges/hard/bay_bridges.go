package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type bridge struct {
	id             int
	x1, y1, x2, y2 float64
	inter          []int
}

type ascend []bridge

func (s ascend) Len() int { return len(s) }
func (s ascend) Less(i, j int) bool {
	return s[i].id < s[j].id
}
func (s ascend) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func intersect(a bridge, b bridge) bool {
	sx1, sy1 := a.x2-a.x1, a.y2-a.y1
	sx2, sy2 := b.x2-b.x1, b.y2-b.y1
	sa := (sx1*(a.y1-b.y1) - sy1*(a.x1-b.x1)) / (sx1*sy2 - sx2*sy1)
	sb := (sx2*(a.y1-b.y1) - sy2*(a.x1-b.x1)) / (sx1*sy2 - sx2*sy1)
	return sa >= 0 && sa <= 1 && sb >= 0 && sb <= 1
}

var maxlen int

func bay(a []bridge, b []int, c int) []int {
	if c == len(a) {
		if len(b) > maxlen {
			maxlen = len(b)
		}
		return b
	}

	var res1 []int
	if maxlen < len(b)+len(a)-c {
		res1 = bay(a, b, c+1)
		for i := 0; i < len(a[c].inter); i++ {
			if ix := sort.SearchInts(b, a[c].inter[i]); ix < len(b) && b[ix] == a[c].inter[i] {
				return res1
			}
		}
		res2 := bay(a, append(b, a[c].id), c+1)
		if len(res2) > len(res1) {
			return res2
		}
	}
	return res1
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	var bridges []bridge
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var nb bridge
		fmt.Sscanf(scanner.Text(), "%d: ([%f,%f], [%f,%f])", &nb.id, &nb.x1, &nb.y1, &nb.x2, &nb.y2)
		bridges = append(bridges, nb)
	}
	sort.Sort(ascend(bridges))
	for i := 0; i < len(bridges)-1; i++ {
		for j := i + 1; j < len(bridges); j++ {
			if intersect(bridges[i], bridges[j]) {
				bridges[i].inter = append(bridges[i].inter, bridges[j].id)
				bridges[j].inter = append(bridges[j].inter, bridges[i].id)
			}
		}
	}

	m := bay(bridges, []int{}, 0)
	for _, i := range m {
		fmt.Println(i)
	}
}
