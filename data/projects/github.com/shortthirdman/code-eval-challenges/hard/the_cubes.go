package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type seg struct {
	pos int
	to  []int
}

var floor, holes map[int]bool

func contains(a []seg, b int) bool {
	for _, i := range a {
		if i.pos == b {
			return true
		}
	}
	return false
}

func maketo(p, n int) (to []int) {
	p2 := p % (n * n)
	if p2/n < n-1 && floor[p+n] {
		to = append(to, p+n)
	}
	if p2/n > 0 && floor[p-n] {
		to = append(to, p-n)
	}
	if p2%n < n-1 && floor[p+1] {
		to = append(to, p+1)
	}
	if p2%n > 0 && floor[p-1] {
		to = append(to, p-1)
	}
	if p/(n*n) < n-1 && holes[p+n*n] {
		to = append(to, p+n*n)
	}
	if p/(n*n) > 0 && holes[p] {
		to = append(to, p-n*n)
	}
	return to
}

func main() {
	var n int
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()

	scanner := bufio.NewScanner(data)
	for scanner.Scan() {
		var (
			downside, upside []seg
			r                int
		)
		s := strings.Split(scanner.Text(), ";")
		fmt.Sscan(s[0], &n)
		floor, holes = make(map[int]bool, n*n*n), make(map[int]bool, n*n*n)

		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				x := (i*n + j) * n
				for k := 0; k < n; k++ {
					switch s[1][x+k] {
					case 'o':
						holes[x+k] = true
						fallthrough
					case ' ':
						floor[x+k] = true
						if j == 0 || j == n-1 || k == 0 || k == n-1 {
							if i == 0 {
								downside = append(downside, seg{pos: x + k})
							} else if i == n-1 {
								upside = append(upside, seg{pos: x + k})
							}
						}
					}
				}
			}
		}
		for ix, i := range downside {
			downside[ix].to = maketo(i.pos, n)
			floor[i.pos] = false
		}
		for ix, i := range upside {
			upside[ix].to = maketo(i.pos, n)
			floor[i.pos] = false
		}
		for x := 2; true; x += 2 {
			if len(downside) == 0 {
				r = 0
				break
			}
			var newdown []seg
			for _, i := range downside {
				for _, j := range i.to {
					if contains(upside, j) {
						r = x
						break
					}
					if floor[j] {
						newdown = append(newdown, seg{j, maketo(j, n)})
						floor[j] = false
					}
				}
				if r > 0 {
					break
				}
			}
			if r > 0 {
				break
			}
			downside = newdown

			if len(upside) == 0 {
				r = 0
				break
			}
			var newup []seg
			for _, i := range upside {
				for _, j := range i.to {
					if contains(downside, j) {
						r = x + 1
						break
					}
					if floor[j] {
						newup = append(newup, seg{j, maketo(j, n)})
						floor[j] = false
					}
				}
				if r > 0 {
					break
				}
			}
			if r > 0 {
				break
			}
			upside = newup
		}
		fmt.Println(r)
	}
}
