package main

import (
	"strings"
	"testing"
)

func TestStepwise(t *testing.T) {
	for k, v := range map[string]string{
		"cat dog hello":        "h *e **l ***l ****o",
		"stop football play":   "f *o **o ***t ****b *****a ******l *******l",
		"music is my life":     "m *u **s ***i ****c",
		"asdf 0123456789 qwer": "0 *1 **2 ***3 ****4 *****5 ******6 *******7 ********8 *********9",
		"a b c d e f":          "a"} {
		if r := stepwise(k); r != v {
			t.Errorf("failed: stepwise %s is %s, got %s",
				k, v, r)
		}
	}
}

func stepwise(s string) string {
	var (
		maxw string
		maxl int
	)
	t := strings.Fields(s)
	for _, i := range t {
		if len(i) > maxl {
			maxw, maxl = i, len(i)
		}
	}
	r := make([]string, maxl)
	for i := 0; i < maxl; i++ {
		r[i] = strings.Repeat("*", i) + string(maxw[i])
	}
	return strings.Join(r, " ")
}
