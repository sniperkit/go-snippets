package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestKernighan(t *testing.T) {
	for k, v := range map[int]int{
		0: 0, 1: 1, 2: 1, 13: 3, 42: 3, 127: 7, 8192: 1,
		8200: 2, 4194305: 2} {
		if r := kernighan(k); r != v {
			t.Errorf("failed: kernighan %d is %d, got %d",
				k, v, r)
		}
	}
}

func TestBlackOrWhite(t *testing.T) {
	for k, v := range map[string]string{
		"11 | 11":                   "1x1, 1",
		"1001 | 0110 | 1001 | 0110": "2x2, 2",
		"110 | 101 | 111":           "3x3, 7",
		"000 | 000 | 000":           "1x1, 0",
		"1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111111 | 1111111111111111111111110": "25x25, 624"} {
		if r := blackOrWhite(k); r != v {
			t.Errorf("failed: blackOrWhite %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkKernighan(b *testing.B) {
	for i := 0; i < b.N; i++ {
		kernighan(i)
	}
}

func kernighan(a int) (r int) {
	for ; a > 0; a &= a - 1 {
		r++
	}
	return r
}

func numCar(m []int, n int) (r int) {
	var i uint
	for i = 0; i < uint(n); i++ {
		r += kernighan(m[i] & ((1 << uint(n)) - 1))
	}
	return r
}

func isNum(m []int, n, x, y, nc int) bool {
	var (
		r int
		i uint
	)
	for i = uint(x); i < uint(x+n); i++ {
		r += kernighan(m[i] & (((1 << uint(n)) - 1) << uint(y)))
		if r > nc {
			return false
		}
	}
	return r == nc
}

func blackOrWhite(q string) string {
	var rm, rn int
	s := strings.Split(q, "|")
	m := make([]int, len(s))
	for ix, i := range s {
		fmt.Sscanf(strings.TrimSpace(i), "%b", &m[ix])
	}

	for rm = 1; rm <= len(s); rm++ {
		f := false
		rn = numCar(m, rm)
		for i := 0; i < len(s)-rm+1; i++ {
			for j := 0; j < len(s)-rm+1; j++ {
				if !isNum(m, rm, i, j, rn) {
					f = true
					break
				}
			}
			if f {
				break
			}
		}
		if !f {
			break
		}
	}
	return fmt.Sprintf("%dx%d, %d", rm, rm, rn)
}
