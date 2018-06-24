package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	s, t string
}

func TestIntersect(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"1,2,3,4", "4,5,6"}:               "4",
		tuple{"7,8,9", "8,9,10,11,12"}:          "8,9",
		tuple{"3,4,7", "2,6,12"}:                "",
		tuple{"1,3,5,7,9,11,13", "2,3,5,6,8,9"}: "3,5,9"} {
		if r := intersect(k.s, k.t); r != v {
			t.Errorf("failed: intersect %s;%s is %s, got %s",
				k.s, k.t, v, r)
		}
	}
}

func BenchmarkIntersect(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var (
			s, t []string
			c    uint
		)
		for j := i + 3; j > 0; j /= 4 {
			s = append(s, fmt.Sprint((j&1)<<c))
			t = append(s, fmt.Sprint((j&2)<<c))
			c++
		}
		intersect(strings.Join(s, ","), strings.Join(t, ","))
	}
}

func intersect(s, t string) string {
	var (
		a    int
		x, y []int
		r    []string
	)
	u, v := strings.Split(s, ","), strings.Split(t, ",")
	for _, i := range u {
		fmt.Sscan(i, &a)
		x = append(x, a)
	}
	for _, i := range v {
		fmt.Sscan(i, &a)
		y = append(y, a)
	}
	for i, j := 0, 0; i < len(u) && j < len(v); {
		if x[i] == y[j] {
			r = append(r, fmt.Sprint(x[i]))
			i++
			j++
		} else if x[i] > y[j] {
			j++
		} else {
			i++
		}
	}
	return strings.Join(r, ",")
}
