package main

import (
	"fmt"
	"strings"
	"testing"
)

type tuple struct {
	s, t string
}

func TestFindWriter(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"osSE5Gu0Vi8WRq93UvkYZCjaOKeNJfTyH6tzDQbxFm4M1ndXIPh27wBA rLclpg", " 3 35 27 62 51 27 46 57 26 10 46 63 57 45 15 43 53"}:             "Stephen King 1947",
		tuple{"3Kucdq9bfCEgZGF2nwx8UpzQJyHiOm0hoaYP6ST1WM7Nks5XjrR4IltBeDLV vA", " 2 26 33 55 34 50 33 61 44 28 46 32 28 30 3 50 34 61 40 7 1 31"}: "Kyotaro Nishimura 1930"} {
		if r := findWriter(k.s, k.t); r != v {
			t.Errorf("failed: findWriter %s|%s is %s, got %s",
				k.s, k.t, v, r)
		}
	}
}

func findWriter(s, t string) string {
	var f int
	u := strings.Fields(t)
	r := make([]byte, len(u))
	for ix, i := range u {
		fmt.Sscan(i, &f)
		r[ix] = byte(s[f-1])
	}
	return string(r)
}
