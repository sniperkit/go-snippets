package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type team struct {
	id int
	cs string
}

type teams []team

func (slice teams) Len() int           { return len(slice) }
func (slice teams) Less(i, j int) bool { return slice[i].id < slice[j].id }
func (slice teams) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

func football(p string) string {
	var (
		a int
		r []team
		q []string
	)
	m := make(map[int][]string)
	s := strings.Split(p, "|")
	for ix, i := range s {
		t := strings.Fields(i)
		for _, j := range t {
			fmt.Sscan(j, &a)
			m[a] = append(m[a], fmt.Sprint(ix+1))
		}
	}
	for k, v := range m {
		r = append(r, team{k, strings.Join(v, ",")})
	}
	sort.Sort(teams(r))
	for _, i := range r {
		q = append(q, fmt.Sprintf("%d:%s;", i.id, i.cs))
	}
	return strings.Join(q, " ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		fmt.Println(football(scanner.Text()))
	}
}
