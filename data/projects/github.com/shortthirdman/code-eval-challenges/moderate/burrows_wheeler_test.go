package main

import (
	"sort"
	"testing"
)

func TestBurrowsWheeler(t *testing.T) {
	for k, v := range map[string]string{
		"oooooooo$  ffffffff     ffffffffuuuuuuuuaaaaaaaallllllllbbBbbBBb":                             "Buffalo buffalo Buffalo buffalo buffalo buffalo Buffalo buffalo$",
		"edarddddddddddntensr$  ehhhhhhhhhhhJ aeaaaaaaaaaaalhtf thmbfe           tcwohiahoJ eeec t e ": "James while John had had had had had had had had had had had a better effect on the teacher$",
		"ooooio,io$Nnssshhhjo  ee  o  nnkkkkkkii ":                                                     "Neko no ko koneko, shishi no ko kojishi$"} {
		if r := burrowsWheeler(k); r != v {
			t.Errorf("failed: burrowsWheeler %s is %s, got %s",
				k, v, r)
		}
	}
}

func burrowsWheeler(s string) string {
	e, r := make([]int, len(s)), make([]byte, len(s))
	var k int
	for i := 0; i < len(s); i++ {
		e[i] = int(s[i])<<16 + i
		if s[i] == '$' {
			k = i
		}
	}
	sort.Ints(e)
	for i := 0; i < len(s); i, k = i+1, e[k]%(1<<16) {
		r[i] = byte(e[k] >> 16)
	}
	return string(r)
}
