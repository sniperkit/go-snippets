package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

func contains(a []int, b int) bool {
	ix := sort.SearchInts(a, b)
	return ix < len(a) && a[ix] == b
}

type host struct {
	id   int
	sub  []int
	conn []int
}

type ascend []host

func (s ascend) Len() int { return len(s) }
func (s ascend) Less(i, j int) bool {
	return s[i].id < s[j].id
}
func (s ascend) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type natural [][]int

func (s natural) Len() int { return len(s) }
func (s natural) Less(i, j int) bool {
	for k := 0; ; k++ {
		if s[i][k] != s[j][k] {
			return s[i][k] < s[j][k]
		}
	}
	return false
}
func (s natural) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func a2ip(a1, a2, a3, a4, a5 int) int {
	t := uint(32 - a5)
	return ((a1<<24 + a2<<16 + a3<<8 + a4) >> t) << t
}

func nice(a [][]int) string {
	ret := make([]string, len(a))
	for ix, i := range a {
		r := make([]string, len(i))
		for jx, j := range i {
			r[jx] = fmt.Sprint(j)
		}
		ret[ix] = "[" + strings.Join(r, ", ") + "]"
	}
	return strings.Join(ret, ", ")
}

func main() {
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	var (
		n, a1, a2, a3, a4, a5, src, dst int
		hosts                           []host
	)
	scanner := bufio.NewScanner(data)
	scanner.Split(bufio.ScanWords)
	for l := ""; l != "]}"; {
		var h host
		l = ""
		scanner.Scan()
		fmt.Sscanf(strings.TrimLeft(scanner.Text(), "{"), "%d:", &h.id)
		for l != "]," && l != "]}" {
			scanner.Scan()
			s := strings.TrimLeft(scanner.Text(), "[")
			fmt.Sscanf(s, "'%d.%d.%d.%d/%d'", &a1, &a2, &a3, &a4, &a5)
			h.sub = append(h.sub, a2ip(a1, a2, a3, a4, a5))
			l = s[len(s)-2:]
		}
		sort.Ints(h.sub)
		hosts = append(hosts, h)
		n++
	}
	sort.Sort(ascend(hosts))
	for i := 0; i < n-1; i++ {
		for j := i + 1; j < n; j++ {
			for k := 0; k < len(hosts[i].sub); k++ {
				ix := sort.SearchInts(hosts[j].sub, hosts[i].sub[k])
				if ix < len(hosts[j].sub) && hosts[j].sub[ix] == hosts[i].sub[k] {
					hosts[i].conn = append(hosts[i].conn, hosts[j].id)
					hosts[j].conn = append(hosts[j].conn, hosts[i].id)
				}
			}
		}
		hosts[i].sub = nil
	}
	hosts[n-1].sub = nil
	id2h := make(map[int]int, n)
	for i := 0; i < n; i++ {
		id2h[hosts[i].id] = i
	}

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		for i := 0; i < n; i++ {
			hosts[i].sub = hosts[i].sub[0:0]
		}
		var fin bool
		fmt.Sscanf(scanner.Text(), "%d %d", &src, &dst)
		q := []int{src}
		for len(q) > 0 {
			var r []int
			for _, i := range q {
				for _, j := range hosts[id2h[i]].conn {
					if j == dst {
						fin = true
					}
					if !contains(hosts[id2h[i]].sub, j) && !contains(q, j) {
						hosts[id2h[j]].sub = append(hosts[id2h[j]].sub, i)
						sort.Ints(hosts[id2h[j]].sub)
						if !contains(r, j) {
							r = append(r, j)
							sort.Ints(r)
						}
					}
				}
			}
			if fin {
				break
			}
			q = r
		}
		if !fin {
			fmt.Println("No connection")
			continue
		}
		paths := [][]int{[]int{dst}}
		for i := 0; i < len(paths); i++ {
			for paths[i][0] != src {
				cur := id2h[paths[i][0]]
				for j := 1; j < len(hosts[cur].sub); j++ {
					newpath := make([]int, len(paths[i])+1)
					copy(newpath[1:], paths[i])
					newpath[0] = hosts[cur].sub[j]
					paths = append(paths, newpath)
				}
				paths[i] = append([]int{hosts[cur].sub[0]}, paths[i]...)
			}
		}
		sort.Sort(natural(paths))
		fmt.Println(nice(paths))
	}
}
