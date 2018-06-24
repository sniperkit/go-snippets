package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type team struct {
	id int
	cs string
}

type teams []team

func (slice teams) Len() int           { return len(slice) }
func (slice teams) Less(i, j int) bool { return slice[i].id < slice[j].id }
func (slice teams) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

func TestFootball(t *testing.T) {
	for k, v := range map[string]string{
		"1 2 3 4 | 3 1 | 4 1":         "1:1,2,3; 2:1; 3:1,2; 4:1,3;",
		"19 11 | 19 21 23 | 31 39 29": "11:1; 19:1,2; 21:2; 23:2; 29:3; 31:3; 39:3;"} {
		if r := football(k); r != v {
			t.Errorf("failed: football %s is %s, got %s",
				k, v, r)
		}
	}
}

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
