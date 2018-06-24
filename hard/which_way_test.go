package main

import (
	"math"
	"sort"
	"strings"
	"testing"
)

func TestWhichWay(t *testing.T) {
	for k, v := range map[string]int{
		"**^F | P**P | **** | S***":             11,
		"**^F | P*^P | **** | S***":             12,
		"PS*** | ^^^^* | ***** | *^^^^ | ***FP": 22} {
		if r := whichWay(k); r != v {
			t.Errorf("failed: whichWay %s is %d, got %d",
				k, v, r)
		}
	}
}

type coor struct {
	x, y, d int
	land    bool
}
type coords []coor

func (slice coords) Len() int { return len(slice) }
func (slice coords) Less(i, j int) bool {
	if slice[i].d == slice[j].d {
		return slice[j].land
	}
	return slice[i].d < slice[j].d
}
func (slice coords) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func contains(a []coor, b coor) bool {
	for _, i := range a {
		if b.x == i.x && b.y == i.y && b.land == i.land {
			return true
		}
	}
	return false
}

func whichWay(q string) int {
	s := strings.Split(q, "|")
	for ix, i := range s {
		s[ix] = strings.TrimSpace(i)
	}
	n, l := len(s), -1
	ml, ms := make([][]int, n), make([][]int, n)
	var t, u, p []coor
	for i := 0; i < n; i++ {
		ml[i], ms[i] = make([]int, n), make([]int, n)
		for j := 0; j < n; j++ {
			switch s[i][j] {
			case 'S':
				t = append(t, coor{i, j, 0, true})
				ml[i][j], ms[i][j] = 0, -1
			case 'F':
				u = append(u, coor{i, j, -1, true})
				ml[i][j], ms[i][j] = -1, -1
			case 'P':
				p = append(p, coor{i, j, -1, true})
				p = append(p, coor{i, j, -1, false})
				ml[i][j], ms[i][j] = -1, -1
			case '^':
				ml[i][j], ms[i][j] = math.MaxInt32, math.MaxInt32
			case '*':
				ml[i][j], ms[i][j] = -1, -1
			}
		}
	}
	for len(t) > 0 {
		cc := t[0]
		t = t[1:]
		if contains(u, cc) {
			l = ml[cc.x][cc.y]
			break
		} else if contains(p, cc) {
			if cc.land {
				if ms[cc.x][cc.y] < 0 {
					ms[cc.x][cc.y] = ml[cc.x][cc.y] + 1
					t = append(t, coor{cc.x, cc.y, cc.d + 1, false})
				}
			} else {
				if ml[cc.x][cc.y] < 0 {
					ml[cc.x][cc.y] = ms[cc.x][cc.y] + 1
					t = append(t, coor{cc.x, cc.y, cc.d + 1, true})
				}
			}
		}
		if cc.land {
			if cc.x > 0 && ml[cc.x-1][cc.y] < 0 {
				ml[cc.x-1][cc.y] = ml[cc.x][cc.y] + 2
				t = append(t, coor{cc.x - 1, cc.y, cc.d + 2, true})
			}
			if cc.x < n-1 && ml[cc.x+1][cc.y] < 0 {
				ml[cc.x+1][cc.y] = ml[cc.x][cc.y] + 2
				t = append(t, coor{cc.x + 1, cc.y, cc.d + 2, true})
			}
			if cc.y > 0 && ml[cc.x][cc.y-1] < 0 {
				ml[cc.x][cc.y-1] = ml[cc.x][cc.y] + 2
				t = append(t, coor{cc.x, cc.y - 1, cc.d + 2, true})
			}
			if cc.y < n-1 && ml[cc.x][cc.y+1] < 0 {
				ml[cc.x][cc.y+1] = ml[cc.x][cc.y] + 2
				t = append(t, coor{cc.x, cc.y + 1, cc.d + 2, true})
			}
		} else {
			if cc.x > 0 && ms[cc.x-1][cc.y] < 0 {
				ms[cc.x-1][cc.y] = ms[cc.x][cc.y] + 1
				t = append(t, coor{cc.x - 1, cc.y, cc.d + 1, false})
			}
			if cc.x < n-1 && ms[cc.x+1][cc.y] < 0 {
				ms[cc.x+1][cc.y] = ms[cc.x][cc.y] + 1
				t = append(t, coor{cc.x + 1, cc.y, cc.d + 1, false})
			}
			if cc.y > 0 && ms[cc.x][cc.y-1] < 0 {
				ms[cc.x][cc.y-1] = ms[cc.x][cc.y] + 1
				t = append(t, coor{cc.x, cc.y - 1, cc.d + 1, false})
			}
			if cc.y < n-1 && ms[cc.x][cc.y+1] < 0 {
				ms[cc.x][cc.y+1] = ms[cc.x][cc.y] + 1
				t = append(t, coor{cc.x, cc.y + 1, cc.d + 1, false})
			}
		}
		sort.Sort(coords(t))
	}
	return l
}
