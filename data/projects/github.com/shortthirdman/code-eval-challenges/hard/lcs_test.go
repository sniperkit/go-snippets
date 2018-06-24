package main

import "testing"

type tuple struct {
	p, q string
}

func TestLcs(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"XMJYAUZ", "MZJAWXU"}:                                                           "MJAU",
		tuple{"ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890", "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"}: "ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890",
		tuple{"ATEST", "TEISTG"}:                                                              "TEST"} {
		if r := lcs(k.p, k.q); r != v {
			t.Errorf("failed: lcs %s;%s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

func bt(c [][]int, x, y string, i, j int) string {
	if i < 1 || j < 1 {
		return ""
	}
	if x[i-1] == y[j-1] {
		return bt(c, x, y, i-1, j-1) + string(x[i-1])
	}
	if c[i][j-1] > c[i-1][j] {
		return bt(c, x, y, i, j-1)
	}
	return bt(c, x, y, i-1, j)
}

func max(i, j int) int {
	if i > j {
		return i
	}
	return j
}

func lcs(p, q string) string {
	c := make([][]int, len(p)+1)
	c[0] = make([]int, len(q)+1)
	for ix, i := range p {
		c[ix+1] = make([]int, len(q)+1)
		for jx, j := range q {
			if i == j {
				c[ix+1][jx+1] = c[ix][jx] + 1
			} else {
				c[ix+1][jx+1] = max(c[ix+1][jx], c[ix][jx+1])
			}
		}
	}
	return bt(c, p, q, len(p), len(q))
}
