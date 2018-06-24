package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

func TestNumberOperations(t *testing.T) {
	for k, v := range map[string]bool{
		"44 6 1 49 47":   false,
		"34 1 49 2 21":   true,
		"31 38 27 51 18": false} {
		if r := numberOperations(k); r != v {
			t.Errorf("failed: numberOperations %s is %t, got %t",
				k, v, r)
		}
	}
}

func nextPermu(s string) (string, bool) {
	k, l := -1, 0
	for i := len(s) - 2; i >= 0; i-- {
		if s[i] < s[i+1] {
			k = i
			break
		}
	}
	if k == -1 {
		return s, true
	}
	for i := len(s) - 1; i > 0; i-- {
		if s[i] > s[k] {
			l = i
			break
		}
	}
	s = strings.Join([]string{s[0:k], string(s[l]), s[k+1 : l], string(s[k]), s[l+1 : len(s)]}, "")
	l = len(s) - k - 1
	for i := 0; l-1-i > i; i++ {
		s = strings.Join([]string{s[0 : k+1+i], string(s[len(s)-1-i]), s[k+2+i : len(s)-1-i], string(s[k+1+i]), s[len(s)-i : len(s)]}, "")
	}
	return s, false
}

func numberOperations(q string) bool {
	s := strings.Fields(q)
	t := make([]int, len(s))
	var u string
	for i := range s {
		fmt.Sscan(s[i], &t[i])
	}
	sort.Ints(t)
	for i := range s {
		u += fmt.Sprintf("%c", t[i])
	}
	v := u
	var f, g bool
	for !f && !g {
		for o1 := 0; o1 < 3 && !g; o1++ {
			r1 := int(v[0])
			switch o1 {
			case 0:
				r1 += int(v[1])
			case 1:
				r1 -= int(v[1])
			case 2:
				r1 *= int(v[1])
			}
			for o2 := 0; o2 < 3 && !g; o2++ {
				r2 := r1
				switch o2 {
				case 0:
					r2 += int(v[2])
				case 1:
					r2 -= int(v[2])
				case 2:
					r2 *= int(v[2])
				}
				for o3 := 0; o3 < 3 && !g; o3++ {
					r3 := r2
					switch o3 {
					case 0:
						r3 += int(v[3])
					case 1:
						r3 -= int(v[3])
					case 2:
						r3 *= int(v[3])
					}
					for o4 := 0; o4 < 3; o4++ {
						r4 := r3
						switch o4 {
						case 0:
							r4 += int(v[4])
						case 1:
							r4 -= int(v[4])
						case 2:
							r4 *= int(v[4])
						}
						if r4 == 42 {
							g = true
							break
						}
					}
				}
			}
		}
		v, f = nextPermu(v)
	}
	return g
}
