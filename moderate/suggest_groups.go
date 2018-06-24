package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type user struct {
	id      int
	friends []int
	groups  []int
}

type byid []user

func (s byid) Len() int { return len(s) }
func (s byid) Less(i, j int) bool {
	return s[i].id < s[j].id
}
func (s byid) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func contains(a []int, b int) bool {
	ix := sort.SearchInts(a, b)
	return ix < len(a) && a[ix] == b
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	var (
		userlist            []user
		gnames, unames, ret []string
	)
	gid, uid := map[string]int{}, map[string]int{}
	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		s := strings.Split(scanner.Text(), ":")
		var cuid int
		if _, f := uid[s[0]]; f {
			cuid = uid[s[0]]
		} else {
			cuid = len(unames)
			uid[s[0]], unames = cuid, append(unames, s[0])
		}
		t, u := strings.Split(s[1], ","), strings.Split(s[2], ",")
		newuser := user{cuid, make([]int, len(t)), make([]int, len(u))}
		for ix, i := range t {
			if _, f := uid[i]; !f {
				uid[i], unames = len(unames), append(unames, i)
			}
			newuser.friends[ix] = uid[i]
		}
		for ix, i := range u {
			if _, f := gid[i]; !f {
				gid[i], gnames = len(gnames), append(gnames, i)
			}
			newuser.groups[ix] = gid[i]
		}
		sort.Ints(newuser.groups)
		userlist = append(userlist, newuser)
	}
	sort.Sort(byid(userlist))
	for _, i := range userlist {
		var suggest []string
		for j := range gnames {
			if !contains(i.groups, j) {
				var ing int
				for _, k := range i.friends {
					if contains(userlist[k].groups, j) {
						ing++
						if 2*ing >= len(i.friends) {
							suggest = append(suggest, gnames[j])
							break
						}
					}
				}
			}
		}
		if len(suggest) > 0 {
			sort.Strings(suggest)
			ret = append(ret, unames[i.id]+":"+strings.Join(suggest, ","))
		}
	}
	sort.Strings(ret)
	fmt.Println(strings.Join(ret, "\n"))
}
