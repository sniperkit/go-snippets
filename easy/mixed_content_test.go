package main

import (
	"strconv"
	"strings"
	"testing"
)

func TestMixedContent(t *testing.T) {
	for k, v := range map[string]string{
		"8,33,21,0,16,50,37,0,melon,7,apricot,peach,pineapple,17,21": "melon,apricot,peach,pineapple|8,33,21,0,16,50,37,0,7,17,21",
		"24,13,14,43,41":                                             "24,13,14,43,41",
		"asd,sl33p,qwer,3e4,yxcv":                                    "asd,sl33p,qwer,3e4,yxcv"} {
		if r := mixedContent(k); r != v {
			t.Errorf("failed: mixedContent %s is %s, got %s",
				k, v, r)
		}
	}
}

func mixedContent(q string) (r string) {
	var t, r1, r2 []string
	t = strings.Split(q, ",")
	for _, i := range t {
		_, err := strconv.ParseUint(i, 10, 64)
		if err != nil {
			r1 = append(r1, i)
		} else {
			r2 = append(r2, i)
		}
	}
	if len(r1) > 0 {
		r = strings.Join(r1, ",")
		if len(r2) > 0 {
			r += "|"
		}
	}
	if len(r2) > 0 {
		r += strings.Join(r2, ",")
	}
	return r
}
