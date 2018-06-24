package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	a, n int
}

func TestDigitStatistics(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{2, 5}:             "0: 0, 1: 0, 2: 2, 3: 0, 4: 1, 5: 0, 6: 1, 7: 0, 8: 1, 9: 0",
		tuple{3, 1000000000000}: "0: 0, 1: 250000000000, 2: 0, 3: 250000000000, 4: 0, 5: 0, 6: 0, 7: 250000000000, 8: 0, 9: 250000000000",
		tuple{4, 999999999997}:  "0: 0, 1: 0, 2: 0, 3: 0, 4: 499999999999, 5: 0, 6: 499999999998, 7: 0, 8: 0, 9: 0",
		tuple{5, 347}:           "0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 347, 6: 0, 7: 0, 8: 0, 9: 0",
		tuple{6, 347}:           "0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 347, 7: 0, 8: 0, 9: 0",
		tuple{7, 123}:           "0: 0, 1: 30, 2: 0, 3: 31, 4: 0, 5: 0, 6: 0, 7: 31, 8: 0, 9: 31",
		tuple{8, 6}:             "0: 0, 1: 0, 2: 1, 3: 0, 4: 2, 5: 0, 6: 1, 7: 0, 8: 2, 9: 0",
		tuple{9, 1}:             "0: 0, 1: 0, 2: 0, 3: 0, 4: 0, 5: 0, 6: 0, 7: 0, 8: 0, 9: 1"} {
		if r := digitStatistics(k.a, k.n); r != v {
			t.Errorf("failed: digitStatistics %d %d is %s, got %s",
				k.a, k.n, v, r)
		}
	}
}

var stats [][]int

func init() {
	stats = [][]int{[]int{2, 4, 8, 6}, []int{3, 9, 7, 1},
		[]int{4, 6}, []int{5}, []int{6},
		[]int{7, 9, 3, 1}, []int{8, 4, 2, 6},
		[]int{9, 1}}
}

func digitStatistics(a, n int) string {
	a -= 2
	m := n / len(stats[a])
	q, r := make([]int, 10), make([]string, 10)
	for _, i := range stats[a] {
		q[i] += m
	}
	for i := 0; i < n%len(stats[a]); i++ {
		q[stats[a][i]]++
	}
	for i := 0; i < 10; i++ {
		r[i] = fmt.Sprintf("%d: %d", i, q[i])
	}
	return strings.Join(r, ", ")
}
