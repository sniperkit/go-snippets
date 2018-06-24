package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestStrsubs(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"10011011001", "0110,1001,1001,0,10,11"}:   "11100110",
		tuple{"1001101100010", "0110,1001,1001,0,10,11"}: "1110011110011"} {
		if r := strsubs(k.p, k.q); r != v {
			t.Errorf("failed: strsubs %s %s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

type segment struct {
	s string
	r bool
}

func strsubs(p, q string) (r string) {
	t := strings.Split(q, ",")
	trans := []segment{segment{p, false}}
	for len(t) > 0 {
		for i := len(trans) - 1; i >= 0; i-- {
			if trans[i].r == false && strings.Contains(trans[i].s, t[0]) {
				u := strings.Split(trans[i].s, t[0])
				var v []segment
				for jx, j := range u {
					if len(j) > 0 {
						v = append(v, segment{j, false})
					}
					if jx != len(u)-1 {
						v = append(v, segment{t[1], true})
					}
				}
				if i > 0 {
					if i < len(trans)-1 {
						trans = append(trans[0:i], append(v, trans[i+1:len(trans)]...)...)
					} else {
						trans = append(trans[0:i], v...)
					}
				} else {
					if i < len(trans)-1 {
						trans = append(v, trans[i+1:len(trans)]...)
					} else {
						trans = v
					}
				}
			}
		}
		t = t[2:len(t)]
	}
	for _, i := range trans {
		r += i.s
	}
	return r
}
