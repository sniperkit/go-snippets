package main

import (
	"strings"
	"testing"
)

type tuple struct {
	a, b int
	q    string
}

func TestSpiral(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{3, 3, "1 2 3 4 5 6 7 8 9"}:                                                                                "1 2 3 6 9 8 7 4 5",
		tuple{5, 7, "1 2 3 4 5 6 7 8 9 10 11 12 13 14 15 16 17 18 asd 20 21 22 23 24 25 26 27 28 29 30 31 32 33 34 35"}: "1 2 3 4 5 6 7 14 21 28 35 34 33 32 31 30 29 22 15 8 9 10 11 12 13 20 27 26 25 24 23 16 17 18 asd",
		tuple{2, 5, "1 2 3 4 5 6 7 8 9 10"}:                                                                             "1 2 3 4 5 10 9 8 7 6",
		tuple{5, 2, "1 2 3 4 5 6 7 8 9 10"}:                                                                             "1 2 4 6 8 10 9 7 5 3",
		tuple{1, 1, "1"}:                                                                                                "1"} {
		if r := spiral(k.a, k.b, k.q); r != v {
			t.Errorf("failed: spiral %d;%d;%s is %s, got %s",
				k.a, k.b, k.q, v, r)
		}
	}
}

func spiral(a, b int, q string) string {
	var (
		r            []string
		i, j, tn, tw int
	)
	t := strings.Fields(q)
	te, ts := b-1, a-1
	for {
		for j <= te {
			r = append(r, t[i*b+j])
			j++
		}
		j--
		i++
		tn++
		if i > ts {
			break
		}
		for i <= ts {
			r = append(r, t[i*b+j])
			i++
		}
		i--
		j--
		te--
		if j < tw {
			break
		}
		for j >= tw {
			r = append(r, t[i*b+j])
			j--
		}
		j++
		i--
		ts--
		if i < tn {
			break
		}
		for i >= tn {
			r = append(r, t[i*b+j])
			i--
		}
		i++
		j++
		tw++
		if j > te {
			break
		}
	}
	return strings.Join(r, " ")
}
