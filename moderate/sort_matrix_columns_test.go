package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type column []int
type columns []column

func (slice columns) Len() int { return len(slice) }
func (slice columns) Less(i, j int) bool {
	var k int
	for k < len(slice[i])-1 && slice[i][k] == slice[j][k] {
		k++
	}
	return slice[i][k] < slice[j][k]
}
func (slice columns) Swap(i, j int) { slice[i], slice[j] = slice[j], slice[i] }

func TestSortMatrix(t *testing.T) {
	for k, v := range map[string]string{
		"-3 29 -3 | -17 69 -17 | 44 3 8":                               "-3 -3 29 | -17 -17 69 | 8 44 3",
		"25 39 -26 -21 | -81 -98 -91 27 | 32 -87 67 98 | -90 -79 18 9": "-26 -21 25 39 | -91 27 -81 -98 | 67 98 32 -87 | 18 9 -90 -79",
		"26 -10 39 | -62 66 97 | 22 85 36":                             "-10 26 39 | 66 -62 97 | 85 22 36"} {
		if r := sortMatrix(k); r != v {
			t.Errorf("failed: sortMatrix %s is %s, got %s",
				k, v, r)
		}
	}
}

func sortMatrix(q string) string {
	var i, j int
	s := strings.Split(q, " | ")
	n := strings.Count(s[0], " ") + 1
	m := make([]column, n)
	for i = 0; i < n; i++ {
		m[i] = make(column, n)
	}
	for i = 0; i < n; i++ {
		t := strings.Fields(s[i])
		for j = 0; j < n; j++ {
			fmt.Sscan(t[j], &m[j][i])
		}
	}
	sort.Sort(columns(m))
	r, t := make([]string, n), make([]string, n)
	for i = 0; i < n; i++ {
		for j = 0; j < n; j++ {
			t[j] = fmt.Sprint(m[j][i])
		}
		r[i] = strings.Join(t, " ")
	}
	return strings.Join(r, " | ")
}
