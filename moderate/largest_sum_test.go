package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestLargestSum(t *testing.T) {
	for k, v := range map[string]int{
		"-10, 2, 3, -2, 0, 5, -15": 8,
		"2,3,-2,-1,10":             12,
		"1,2,3,-8,2,2,3":           7,
		"-2,-2,-2,-3,-2":           -2} {
		if r := largestSum(k); r != v {
			t.Errorf("failed: largestSum %s is %d, got %d",
				k, v, r)
		}
	}
}

func largestSum(q string) int {
	s := strings.Split(q, ",")
	var l, m, v int
	fmt.Sscan(strings.TrimSpace(s[0]), &m)
	for _, i := range s {
		fmt.Sscan(strings.TrimSpace(i), &v)
		if v+l > m {
			m = v + l
		}
		if v+l > 0 {
			l += v
		} else {
			l = 0
		}
	}
	return m
}
