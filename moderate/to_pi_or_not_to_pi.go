package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
)

type pwi struct {
	n, nx int
}
type pwis []pwi

func (slice pwis) Len() int           { return len(slice) }
func (slice pwis) Less(i, j int) bool { return slice[i].n < slice[j].n }
func (slice pwis) Swap(i, j int)      { slice[i], slice[j] = slice[j], slice[i] }

func main() {
	n := make([]pwi, 30)
	data, err := os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer data.Close()
	scanner := bufio.NewScanner(data)
	var i, mn int
	for scanner.Scan() {
		var x int
		fmt.Sscanf(scanner.Text(), "%d", &x)
		n[i] = pwi{x, i}
		if mn < x {
			mn = x
		}
		i++
	}
	sort.Sort(pwis(n))
	var (
		l, s    int
		j, p, m int
		a       []int
		r       []string
		x, q, k int
	)
	l = 10*mn/3 + 1
	a = make([]int, l)
	r = make([]string, 30)
	for j = mn; j > 0; {
		q = 0
		k = l + l - 1
		for i = l; i > 0; i-- {
			if j == mn {
				x = 20 + q*i
			} else {
				x = 10*a[i-1] + q*i
			}
			q = x / k
			a[i-1] = x - q*k
			k -= 2
		}
		k = x % 10
		if k == 9 {
			m++
		} else {
			if j > 0 {
				for s < 30 && (mn-j+1 == n[s].n) {
					r[n[s].nx] = fmt.Sprintf("%d", p+x/10)
					s++
				}
				j--
			}
			for ; m > 0; m-- {
				if j > 0 {
					for s < 30 && (mn-j+1 == n[s].n) {
						if x >= 10 {
							r[n[s].nx] = "0"
						} else {
							r[n[s].nx] = "9"
						}
						s++
					}
					j--
				}
			}
			p = k
		}
	}
	fmt.Println(strings.Join(r, "\n"))
}
