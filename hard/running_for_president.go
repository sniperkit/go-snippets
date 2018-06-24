package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type state struct {
	votes int
	value map[int]float64 // issue.id -> fraction of votes
}

type issue struct {
	id   int
	name string
	cost int
	frac float64 // total fractional votes
}

type descend []issue

func (s descend) Len() int { return len(s) }
func (s descend) Less(i, j int) bool {
	if s[i].frac == s[j].frac {
		return s[i].cost < s[j].cost
	}
	return s[i].frac > s[j].frac
}
func (s descend) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

var (
	states []state
	issues []issue
	id2ix  map[int]int
)

func init() {
	id2ix = make(map[int]int)
}

func popular(is []int, s state) bool {
	var v float64
	for _, i := range is {
		if _, f := s.value[i]; f {
			v += s.value[i]
		}
	}
	if v >= 0.51 {
		return true
	}
	return false
}

func elect(is []int) (votes, cost int) {
	for _, i := range is {
		cost += issues[id2ix[i]].cost
	}
	for _, i := range states {
		if popular(is, i) {
			votes += i.votes
		}
	}
	return votes, cost
}

type task struct {
	level  int
	votes  int
	cost   int
	issues []int
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	var niss, tvot int
	scanner.Scan()
	fmt.Sscanf(scanner.Text(), "Social issues: %d", &niss)
	issues = make([]issue, niss)
	iss := make(map[string]int)
	for i := 0; i < niss; i++ {
		scanner.Scan()
		for scanner.Text() == "" {
			scanner.Scan()
		}
		var cost int
		s := strings.Split(scanner.Text(), ": ")
		fmt.Sscanf(s[1], "%d", &cost)
		issues[i] = issue{i, strings.TrimRight(s[0], ":"), cost, 0}
		iss[issues[i].name] = i
	}
	for scanner.Scan() {
		name := scanner.Text()
		if len(name) == 0 {
			continue
		}
		scanner.Scan()
		var votes int
		fmt.Sscanf(scanner.Text(), "Votes: %d", &votes)
		var total float64
		value := make(map[int]float64)
		for scanner.Scan() {
			if scanner.Text() == "" {
				break
			}
			s := strings.Split(scanner.Text(), ": ")
			var v float64
			fmt.Sscanf(s[1], "%f", &v)
			total += v
			value[iss[s[0]]] = v
		}
		for k, v := range value {
			if v > 0 {
				value[k] = v / total
				issues[k].frac += value[k] * float64(votes)
			}
		}
		states = append(states, state{votes, value})
		tvot += votes
	}
	sort.Sort(descend(issues))
	for ix, i := range issues {
		id2ix[i.id] = ix
	}

	votes, cost := elect([]int{issues[0].id})
	todo := []task{task{}, task{0, votes, cost, []int{issues[0].id}}}
	best, bestc, bestl := niss, 0, make([]int, niss)
	for i := 0; i < niss; i++ {
		bestc += issues[i].cost
		bestl[i] = i
	}
	for len(todo) > 0 {
		curr := todo[len(todo)-1]
		todo = todo[:len(todo)-1]
		if len(curr.issues) > best {
			continue
		}
		if curr.level == niss-2 {
			// finish
			curr.level++
			votes, cost = elect(append(curr.issues, issues[curr.level].id))
			if votes > tvot/2 {
				if len(curr.issues)+1 < best || (len(curr.issues)+1 == best && cost < bestc) {
					best, bestc, bestl = len(curr.issues)+1, cost, append(curr.issues, issues[curr.level].id)
				}
			}
		} else {
			// advance
			curr.level++
			todo = append(todo, curr)
			votes, cost = elect(append(curr.issues, issues[curr.level].id))
			if votes > tvot/2 {
				if len(curr.issues)+1 < best || (len(curr.issues)+1 == best && cost < bestc) {
					best, bestc, bestl = len(curr.issues)+1, cost, append(curr.issues, issues[curr.level].id)
				}
			} else {
				todo = append(todo, task{curr.level, votes, cost, append(curr.issues, issues[curr.level].id)})
			}
		}
	}
	var ret []string
	for _, i := range bestl {
		ret = append(ret, issues[id2ix[i]].name)
	}
	sort.Strings(ret)
	fmt.Println(strings.Join(ret, "\n"))
}
