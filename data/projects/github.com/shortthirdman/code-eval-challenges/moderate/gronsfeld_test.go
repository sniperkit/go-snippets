package main

import (
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestGronsfeld(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"31415", "HYEMYDUMPS"}:         "EXALTATION",
		tuple{"45162", `M%muxi%dncpqftiix"`}: "I love challenges!",
		tuple{"14586214", "Uix!&kotvx3"}:     "Test input."} {
		if r := gronsfeld(k.p, k.q); r != v {
			t.Errorf("failed: gronsfeld %s;%s is %s, got %s",
				k.p, k.q, v, r)
		}
	}
}

const k = ` !"#$%&'()*+,-./0123456789:<=>?@ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz`

func gronsfeld(p, q string) string {
	r := make([]byte, len(q))
	for i := 0; i < len(q); i++ {
		r[i] = k[(len(k)+strings.IndexRune(k, rune(q[i]))-int(p[i%len(p)]-'0'))%len(k)]
	}
	return string(r)
}
