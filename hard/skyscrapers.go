package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type house struct {
	x, h, y int
}
type houses []house

func (s houses) Len() int { return len(s) }
func (s houses) Less(i, j int) bool {
	if s[i].x == s[j].x {
		if s[i].h == s[j].h {
			return s[i].y > s[j].y
		}
		return s[i].h > s[j].h
	}
	return s[i].x < s[j].x
}
func (s houses) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type houseHeap []house

func (h houseHeap) Len() int           { return len(h) }
func (h houseHeap) Less(i, j int) bool { return h[i].y < h[j].y }
func (h houseHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *houseHeap) Push(x interface{}) {
	*h = append(*h, x.(house))
}
func (h *houseHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h houseHeap) maxremain(y int) (ret int) {
	for _, i := range h {
		if i.y > y && i.h > ret {
			ret = i.h
		}
	}
	return ret
}

func main() {
	var x, h, y int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(strings.Replace(scanner.Text(), " ", "", -1), ";")

		t := make([]house, len(s))
		for ix, i := range s {
			fmt.Sscanf(i, "(%d,%d,%d)", &x, &h, &y)
			t[ix] = house{x, h, y}
		}
		sort.Sort(houses(t))
		var res []string
		ch, cx, curr, ly := 0, 0, &houseHeap{}, 0
		heap.Init(curr)
		for _, i := range t {
			for curr.Len() > 0 && (*curr)[0].y < i.x {
				if ly < (*curr)[0].y {
					p := heap.Pop(curr).(house)
					nh := curr.maxremain(p.y)
					if ch > nh {
						ly = p.y
						ch = nh
						res = append(res, fmt.Sprint(ly))
						res = append(res, fmt.Sprint(ch))
					}
				}
				for curr.Len() > 0 && ly == (*curr)[0].y {
					heap.Pop(curr)
				}
			}
			if ch < i.h {
				ch, cx = i.h, i.x
				res = append(res, fmt.Sprint(cx))
				res = append(res, fmt.Sprint(ch))
				ly = cx
			}
			heap.Push(curr, i)
		}
		for curr.Len() > 0 {
			if ly < (*curr)[0].y {
				p := heap.Pop(curr).(house)
				nh := curr.maxremain(p.y)
				if ch > nh {
					ly = p.y
					ch = nh
					res = append(res, fmt.Sprint(ly))
					res = append(res, fmt.Sprint(ch))
				}
			}
			for curr.Len() > 0 && ly == (*curr)[0].y {
				heap.Pop(curr)
			}
		}
		fmt.Println(strings.Join(res, " "))
	}
}
