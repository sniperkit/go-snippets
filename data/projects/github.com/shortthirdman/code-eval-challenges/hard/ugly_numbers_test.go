package main

import (
	"fmt"
	"testing"
)

func TestPow3(t *testing.T) {
	for k, v := range map[int]int{
		0: 1, 1: 3, 39: 4052555153018976267} {
		if r := pow3(k); r != v {
			t.Errorf("failed: 3^%d = %d, got %d",
				k, v, r)
		}
	}
}

func TestUglyNumbers(t *testing.T) {
	for k, v := range map[string]int{
		"1":             0,
		"9":             1,
		"011":           6,
		"12345":         64,
		"120":           6,
		"10010":         64,
		"123456":        199,
		"9934567890123": 426134} {
		if r := uglyNumbers(k); r != v {
			t.Errorf("failed: uglyNumbers %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkPow3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pow3(i % 40)
	}
}

func BenchmarkPow3Naive(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pow3Naive(i % 40)
	}
}

func BenchmarkUglyNumbers(b *testing.B) {
	for i := 0; i < b.N; i++ {
		uglyNumbers(fmt.Sprint(i))
	}
}

func pow3(x int) int {
	y, ret := 3, 1
	for x > 0 {
		if x&1 > 0 {
			ret *= y
		}
		y *= y
		x >>= 1
	}
	return ret
}

func pow3Naive(a int) int {
	ret := 1
	for i := 0; i < a; i++ {
		ret *= 3
	}
	return ret
}

func ugly(j int, n []int) bool {
	var (
		k  int
		s  int64
		c  = int64(n[0])
		cj = j
		p  = true
	)
	for k < len(n)-1 {
		ops := cj % 3
		cj /= 3
		if ops == 0 {
			c *= 10
		} else {
			if p {
				s += c
				if ops == 1 {
					p = false
				}
			} else {
				s -= c
				if ops == 2 {
					p = true
				}
			}
			c = 0
		}
		k++
		c += int64(n[k])
	}
	if p {
		s += c
	} else {
		s -= c
	}
	return s%2 == 0 || s%3 == 0 || s%5 == 0 || s%7 == 0
}

func uglyNumbers(q string) (r int) {
	t := make([]int, len(q))
	for ix, i := range q {
		t[ix] = int(i - '0')
	}
	p := pow3(len(t) - 1)
	for j := 0; j < p; j++ {
		if ugly(j, t) {
			r++
		}
	}
	return r
}
