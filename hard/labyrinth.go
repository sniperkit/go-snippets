package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type seg struct {
	pos  int
	from *seg
	to   []int
}

func kill(a seg, b int) {
	for ix, i := range a.to {
		if i == b {
			a.to[ix], a.to[len(a.to)-1], a.to = a.to[len(a.to)-1], 0, a.to[:len(a.to)-1]
			break
		}
	}
	if len(a.to) == 0 {
		kill(*a.from, a.pos)
	}
}

func main() {
	var (
		m            map[int]bool
		n            int
		north, south []seg
		path         []int
	)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	for i := 0; scanner.Scan(); i++ {
		if n == 0 {
			n = len(scanner.Text())
			m = make(map[int]bool, n*n)
		}
		for jx, j := range scanner.Text() {
			if j == ' ' {
				m[i*n+jx] = true
			}
		}
	}
	for i := 0; i < n; i++ {
		if m[i] && m[i+n] {
			start := seg{i, nil, []int{i + n}}
			north = append(north, seg{pos: i + n, from: &start})
		}
		if m[(n-1)*n+i] && m[(n-2)*n+i] {
			start := seg{(n-1)*n + i, nil, []int{(n-2)*n + i}}
			south = append(south, seg{pos: (n-2)*n + i, from: &start})
		}
	}
	for {
		nr := len(north)
		for i := 0; i < nr; i++ {
			curr, count := north[i], 0
			checks := []int{curr.pos + n, curr.pos - 1, curr.pos + 1, curr.pos - n}
			for _, check := range checks {
				if check != curr.from.pos && m[check] {
					for _, j := range south {
						if check == j.pos {
							for ; curr.pos > n; curr = *curr.from {
								path = append([]int{curr.pos}, path...)
							}
							path = append([]int{curr.pos}, path...)
							for curr = j; curr.pos < (n-1)*n; curr = *curr.from {
								path = append(path, curr.pos)
							}
							path = append(path, curr.pos)
							goto out
						}
					}
					var f bool
					for jx := i + 1; jx < len(north); jx++ {
						if check == north[jx].pos {
							f = true
							break
						}
					}
					if !f {
						north, curr.to = append(north, seg{pos: check, from: &curr}), append(curr.to, check)
						count++
					}
				}
			}
			if count == 0 {
				kill(*curr.from, curr.pos)
			}
		}
		copy(north, north[nr:])
		north = north[:len(north)-nr]
		sr := len(south)
		for i := 0; i < sr; i++ {
			curr, count := south[i], 0
			checks := []int{curr.pos + n, curr.pos - 1, curr.pos + 1, curr.pos - n}
			for _, check := range checks {
				if check != curr.from.pos && m[check] {
					for _, j := range north {
						if check == j.pos {
							for ; curr.pos < (n-1)*n; curr = *curr.from {
								path = append(path, curr.pos)
							}
							path = append(path, curr.pos)
							for curr = j; curr.pos > n; curr = *curr.from {
								path = append([]int{curr.pos}, path...)
							}
							path = append([]int{curr.pos}, path...)
							goto out
						}
					}
					var f bool
					for jx := i + 1; jx < len(south); jx++ {
						if check == south[jx].pos {
							f = true
							break
						}
					}
					if !f {
						south, curr.to = append(south, seg{pos: check, from: &curr}), append(curr.to, check)
						count++
					}
				}
			}
			if count == 0 {
				kill(*curr.from, curr.pos)
			}
		}
		copy(south, south[sr:])
		south = south[:len(south)-sr]
	}
out:
	sort.Ints(path)
	for count, i := 0, 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if x := i*n + j; m[x] {
				for path[count] < x {
					count++
				}
				if path[count] == x {
					fmt.Print("+")
				} else {
					fmt.Print(" ")
				}
			} else {
				fmt.Print("*")
			}
		}
		fmt.Println()
	}
}
