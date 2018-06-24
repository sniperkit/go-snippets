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

type ascend []visit

func (s ascend) Len() int { return len(s) }
func (s ascend) Less(i, j int) bool {
	return s[i].room < s[j].room
}
func (s ascend) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func d2ts(d string) int {
	var h, m, s int
	fmt.Sscanf(d, "%d:%d:%d", &h, &m, &s)
	return 3600*h + 60*m + s
}

func ts2d(ts int) string {
	return fmt.Sprintf("%02d:%02d:%02d", ts/3600, (ts/60)%60, ts%60)
}

func hallway(room int) int {
	return (room / 100) * 100
}

type visit struct {
	room int
	dura int
}

type agent struct {
	id    byte
	enter int
	time  int
	path  []visit
	curr  int
	index int
}

type prioEventQueue []*agent

func (pq prioEventQueue) Len() int { return len(pq) }
func (pq prioEventQueue) Less(i, j int) bool {
	if pq[i].time == pq[j].time {
		return pq[i].id < pq[j].id
	}
	return pq[i].time < pq[j].time
}
func (pq prioEventQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}
func (pq *prioEventQueue) Push(x interface{}) {
	n := len(*pq)
	a := x.(*agent)
	a.index = n
	*pq = append(*pq, a)
}
func (pq *prioEventQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	a := old[n-1]
	a.index = -1
	*pq = old[0 : n-1]
	return a
}
func (pq *prioEventQueue) update(a *agent, time, curr int, path []visit) {
	a.time = time
	a.path = path
	a.curr = curr
	heap.Fix(pq, a.index)
}

func main() {
	var room, dura int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	// room -1: outside building
	// room  xx00: hallway xth floor
	var (
		staff []agent
		n     int
	)
	hq := make(map[int]int) // map occupied until
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var path []visit
		s := strings.Fields(scanner.Text())
		id, time := s[0][0], d2ts(s[1])
		for i := 2; i < len(s); i += 2 {
			fmt.Sscanf(s[i], "%d", &room)
			hq[room] = 0
			fmt.Sscan(s[i+1], &dura)
			path = append(path, visit{room, dura})
		}
		sort.Sort(ascend(path))
		staff = append(staff, agent{id, time, time + 10, path, 100, n})
		n++
	}
	q := make(prioEventQueue, len(staff))
	for i := range staff {
		q[i] = &staff[i]
	}
	heap.Init(&q)

	for q.Len() > 0 {
		p := heap.Pop(&q).(*agent)
		time := p.time
		if len(p.path) == 0 { // no more rooms to visit
			if p.curr == 100 { // on hallway 1st floor
				p.time += 10
				p.curr = -1 // has left the building
			} else { // in some hallway >1st floor
				p.time += p.curr / 10 // and we're out
				p.curr = -1
			}
		} else {
			if p.curr == hallway(p.path[0].room) {
				if hq[p.path[0].room] > time { // is room occupied?
					p.time = hq[p.path[0].room]
				} else {
					hq[p.path[0].room] = time + p.path[0].dura
					p.time = time + p.path[0].dura + 10
					copy(p.path, p.path[1:])
					p.path = p.path[:len(p.path)-1]
				}
			} else {
				p.time += (p.path[0].room/100)*10 - p.curr/10 + 10
				p.curr = hallway(p.path[0].room)
			}
			heap.Push(&q, p)
		}
	}
	for _, i := range staff {
		fmt.Println(ts2d(i.enter), ts2d(i.time))
	}
}
