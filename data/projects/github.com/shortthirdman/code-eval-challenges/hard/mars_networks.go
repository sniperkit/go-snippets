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

type point struct {
	x, y float64
	id   int
}

func dist(a, b point) float64 {
	x, y := a.x-b.x, a.y-b.y
	return math.Sqrt(x*x + y*y)
}

type pair struct {
	d    float64
	a, b int
}

type ascend []pair

func (s ascend) Len() int { return len(s) }
func (s ascend) Less(i, j int) bool {
	return s[i].d < s[j].d
}
func (s ascend) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func nice(a float64) string {
	return fmt.Sprintf("%.0f", math.Ceil(a))
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Fields(scanner.Text())
		p, seg := make([]point, len(s)), make([]*[]int, len(s))
		for i := range s {
			var x, y float64
			fmt.Sscanf(s[i], "%f,%f", &x, &y)
			p[i], seg[i] = point{x, y, i}, &[]int{i}
		}
		if len(p) < 3 {
			if len(p) < 2 {
				fmt.Println(0)
			} else {
				fmt.Println(nice(dist(p[0], p[1])))
			}
			continue
		}
		dis := make([]pair, len(p)*(len(p)-1)/2)
		for i, c := 0, 0; i < len(p)-1; i++ {
			for j := i + 1; j < len(p); j++ {
				dis[c], c = pair{dist(p[i], p[j]), i, j}, c+1
			}
		}
		sort.Sort(ascend(dis))
		seg[dis[0].a] = &[]int{dis[0].a, dis[0].b}
		seg[dis[0].b] = seg[dis[0].a]
		ret := dis[0].d
		n, c := len(p)-1, 1
		for n > 1 {
			if seg[dis[c].a] != seg[dis[c].b] {
				for _, i := range *seg[dis[c].b] {
					*seg[dis[c].a] = append(*seg[dis[c].a], i)
				}
				temp := make([]int, len(*seg[dis[c].b]))
				copy(temp, *seg[dis[c].b])
				for _, i := range temp {
					seg[i] = seg[dis[c].a]
				}
				ret += dis[c].d
				n--
			}
			c++
		}
		fmt.Println(nice(ret))
	}
}
