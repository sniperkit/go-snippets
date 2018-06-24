package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestReverseGroups(t *testing.T) {
	for k, v := range map[string]string{
		"1,2,3,4,5;2":                  "2,1,4,3,5",
		"1,2,3,4,5;3":                  "3,2,1,4,5",
		"22,23,24,25,26,27,28,29,30;3": "24,23,22,27,26,25,30,29,28"} {
		if r := reverseGroups(k); r != v {
			t.Errorf("failed: reverseGroups %s is %s, got %s",
				k, v, r)
		}
	}
}

func reverseGroups(q string) string {
	var (
		m int
		r []string
	)
	t := strings.Split(q, ";")
	fmt.Sscan(t[1], &m)
	u := strings.Split(t[0], ",")
	mi := m
	for ix := range u {
		var o int
		if mi <= len(u)-ix {
			mi--
			o = mi - ix%m
			if mi == 0 {
				mi = m
			}
		}
		r = append(r, u[ix+o])
	}
	return strings.Join(r, ",")
}
