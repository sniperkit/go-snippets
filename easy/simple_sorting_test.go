package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestSimpleSort(t *testing.T) {
	for k, v := range map[string]string{
		"70.920 -38.797 14.354 99.323 90.374 7.581":   "-38.797 7.581 14.354 70.920 90.374 99.323",
		"-37.507 -3.263 40.079 27.999 65.213 -55.552": "-55.552 -37.507 -3.263 27.999 40.079 65.213",
		"-2.2 2.234 -1.1 -0.3 1.33 0.345":             "-2.200 -1.100 -0.300 0.345 1.330 2.234"} {
		if r := simpleSort(k); r != v {
			t.Errorf("failed: simpleSort %s is %s, got %s",
				k, v, r)
		}
	}
}

func simpleSort(q string) string {
	t := strings.Fields(q)
	u := make([]float64, len(t))
	for ix, i := range t {
		fmt.Sscan(i, &u[ix])
	}
	sort.Float64s(u)
	for ix, i := range u {
		t[ix] = fmt.Sprintf("%.3f", i)
	}
	return strings.Join(t, " ")
}
