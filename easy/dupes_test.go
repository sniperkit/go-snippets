package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestDupes(t *testing.T) {
	for k, v := range map[string]string{
		"1,1,1,2,2,3,3,4,4": "1,2,3,4",
		"3,3,3":             "3",
		"2,3,4,5,5":         "2,3,4,5"} {
		if r := dupes(k); r != v {
			t.Errorf("failed: dupes %s is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkDupes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var t []string
		for j, c := i+1, 1; j > 0; j /= 2 {
			if j%2 > 0 {
				t = append(t, fmt.Sprint(c))
			} else {
				c++
			}
		}
		dupes(strings.Join(t, ","))
	}
}

func dupes(q string) string {
	var (
		p string
		r []string
	)
	t := strings.Split(q, ",")
	for _, i := range t {
		if i != p {
			r = append(r, i)
			p = i
		}
	}
	return strings.Join(r, ",")
}
