package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"sort"
)

type point struct {
	x, y float64
	id   int
}

type horizontal []point

func (s horizontal) Len() int { return len(s) }
func (s horizontal) Less(i, j int) bool {
	return s[i].x < s[j].x
}
func (s horizontal) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type vertical []point

func (s vertical) Len() int { return len(s) }
func (s vertical) Less(i, j int) bool {
	return s[i].y < s[j].y
}
func (s vertical) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func dist(a, b point) float64 {
	x, y := a.x-b.x, a.y-b.y
	return math.Sqrt(x*x + y*y)
}

func bruteForceClosestPair(a []point) (float64, int, int) {
	d, v, w := math.MaxFloat64, -1, -1
	for ix, i := range a {
		for jx, j := range a {
			if ix != jx {
				di := dist(i, j)
				if di < d {
					d, v, w = di, i.id, j.id
				}
			}
		}
	}
	return d, v, w
}

func divideAndConquerClosestPair(h, v []point) (float64, int, int) {
	if len(h) <= 3 {
		return bruteForceClosestPair(h)
	}
	m := len(h) / 2
	mid := (h[m-1].x + h[m].x) / 2
	hl, hr := make([]point, m), make([]point, len(h)-m)
	copy(hl, h[:m+1])
	copy(hr, h[m:])
	var vl, vr []point
	for _, i := range v {
		found := false
		for _, j := range hl {
			if j.id == i.id {
				found = true
				break
			}
		}
		if found {
			vl = append(vl, i)
		} else {
			vr = append(vr, i)
		}
	}
	d, x, y := divideAndConquerClosestPair(hl, vl)
	if d2, x2, y2 := divideAndConquerClosestPair(hr, vr); d2 < d {
		d, x, y = d2, x2, y2
	}
	var s []point
	for _, i := range v {
		if math.Abs(i.x-mid) < d {
			s = append(s, i)
		}
	}
	for i := 0; i < len(s)-1; i++ {
		for j := i + 1; j < i+4 && j < len(s); j++ {
			if d2 := dist(s[i], s[j]); d2 < d {
				d, x, y = d2, s[i].id, s[j].id
			}
		}
	}
	return d, x, y
}

func main() {
	var (
		n    int
		x, y float64
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Sscanf(scanner.Text(), "%d", &n)
		if n == 0 {
			break
		}
		p := make([]point, n)
		for i := 0; i < n; i++ {
			scanner.Scan()
			fmt.Sscanf(scanner.Text(), "%f %f", &x, &y)
			p[i] = point{x, y, i}
		}

		h, v := make([]point, len(p)), make([]point, len(p))
		copy(h, p)
		copy(v, p)
		sort.Sort(horizontal(h))
		sort.Sort(vertical(v))
		dist, _, _ := divideAndConquerClosestPair(h, v)
		if dist < 10000 {
			fmt.Printf("%.4f\n", dist)
		} else {
			fmt.Println("INFINITY")
		}
	}
}
