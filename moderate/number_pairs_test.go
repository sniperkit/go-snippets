package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestNumberPairs(t *testing.T) {
	for k, v := range map[string]string{
		"1,2,3,4,6;5":        "1,4;2,3",
		"2,4,5,6,9,11,15;20": "5,15;9,11",
		"1,2,3,4;50":         "NULL"} {
		if r := numberPairs(k); r != v {
			t.Errorf("failed: numberPairs %s is %s, got %s",
				k, v, r)
		}
	}
}

func numberPairs(q string) string {
	var (
		n int
		r []string
	)
	t := strings.Split(q, ";")
	u := strings.Split(t[0], ",")
	fmt.Sscan(t[1], &n)
	v, vk := make(map[int]bool, len(u)), make([]int, len(u))
	for i := range u {
		fmt.Sscan(u[i], &vk[i])
		v[vk[i]] = true
	}
	for i := 0; i < len(vk) && 2*vk[i] < n; i++ {
		if v[n-vk[i]] {
			r = append(r, fmt.Sprint(vk[i], ",", n-vk[i]))
		}
	}
	if len(r) == 0 {
		return "NULL"
	}
	return strings.Join(r, ";")
}
