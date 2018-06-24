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

type item struct {
	id, value, weight int64
}
type items []item

func (slice items) Len() int { return len(slice) }
func (slice items) Less(i, j int) bool {
	return (float64(slice[i].value) / float64(slice[i].weight)) > (float64(slice[j].value) / float64(slice[j].weight))
}
func (slice items) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

type bäst []int64

func (slice bäst) Len() int           { return len(slice) }
func (slice bäst) Less(i, j int) bool { return slice[i] < slice[j] }
func (slice bäst) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

type task struct {
	level, bound, value, room int64
	stuff                     []int64
}

func boundary(room, value int64, stuff []item, itl []int64) (t, b int64, tl []int64) {
	totalRoom, f := room, false
	t, b = value, value
	for _, i := range stuff {
		if i.weight > totalRoom {
			continue
		} else if i.weight > room {
			if !f {
				b += int64(float64(i.value*room) / float64(i.weight))
				f = true
			}
			continue
		}
		room, t, tl = room-i.weight, t+i.value, append(tl, i.id)
		if !f {
			b += i.value
		}
		if room == 0 {
			break
		}
	}
	return t, b, tl
}

func packageProblem(q string) string {
	var (
		ks, s0, s2 int64
		s1         float64
	)
	ts := strings.Split(q, "(")
	n := len(ts) - 1
	fmt.Sscanf(ts[0], "%d : ", &ks)
	var (
		stuff     []item
		minWeight int64 = math.MaxInt64
	)
	n0 := n
	for i := 1; i <= n0; i++ {
		fmt.Sscanf(ts[i], "%d,%f,$%d)", &s0, &s1, &s2)

		if s1 > float64(ks) {
			n--
		} else {
			stuff = append(stuff, item{s0, s2, int64(s1 * 100)})
		}
		if s2 < minWeight {
			minWeight = s2
		}
	}
	if n == 0 {
		return "-"
	}
	sort.Sort(items(stuff))
	ks *= 100

	t, bound, tl := boundary(ks, int64(0), stuff[1:], []int64{})
	best, bestl := t, make([]int64, len(tl))
	todo := []task{task{bound: bound, room: ks}}
	copy(bestl, tl)

	t, bound, tl = boundary(ks, int64(0), stuff, []int64{})
	if t > best {
		best, bestl = t, make([]int64, len(tl))
		copy(bestl, tl)
	}
	if stuff[0].weight <= ks {
		todo = append(todo, task{int64(0), bound, stuff[0].value, ks - stuff[0].weight, []int64{stuff[0].id}})
	}

	for len(todo) > 0 {
		curr := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		if curr.bound <= best {
			continue
		}

		if curr.level == int64(n-2) {
			if curr.room >= stuff[n-1].weight {
				t = stuff[n-1].value + curr.value
				if t > best {
					best, bestl = t, make([]int64, len(curr.stuff)+1)
					copy(bestl, append(curr.stuff, stuff[n-1].id))
				}
			}
			if curr.value > best {
				best, bestl = curr.value, make([]int64, len(curr.stuff))
				copy(bestl, curr.stuff)
			}
			continue
		}
		t, bound, tl = boundary(curr.room, curr.value, stuff[curr.level+1:], curr.stuff)
		if t > best {
			tl = append(curr.stuff, tl...)
			best, bestl = t, make([]int64, len(tl))
			copy(bestl, tl)
		}
		if bound > best && curr.room > minWeight {
			todo = append(todo, task{curr.level + 1, bound, curr.value, curr.room, curr.stuff})
		}

		if curr.room >= stuff[curr.level+1].weight+minWeight {
			todo = append(todo, task{curr.level + 1, curr.bound, curr.value + stuff[curr.level+1].value, curr.room - stuff[curr.level+1].weight, append(curr.stuff, stuff[curr.level+1].id)})
		}
	}
	var st []string
	sort.Sort(bäst(bestl))
	for i := int64(1); i <= int64(n0); i++ {
		if len(bestl) > 0 && i == bestl[0] {
			st = append(st, fmt.Sprint(i))
			bestl = bestl[1:]
		}
	}
	return strings.Join(st, ",")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		if len(scanner.Text()) > 2 {
			fmt.Println(packageProblem(scanner.Text()))
		}
	}
}
