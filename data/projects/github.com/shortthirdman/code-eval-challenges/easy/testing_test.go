package main

import "testing"

type tuple struct {
	s, t string
}

func TestTest(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{"Heelo Codevval", "Hello Codeeval"}: "Low",
		tuple{"hELLO cODEEVAL", "Hello Codeeval"}: "Critical",
		tuple{"Hello Codeeval", "Hello Codeeval"}: "Done"} {
		if r := test(k.s, k.t); r != v {
			t.Errorf("failed: test %s | %s is %s, got %s",
				k.s, k.t, v, r)
		}
	}
}

func test(s, t string) string {
	var r int
	for i := range s {
		if s[i] != t[i] {
			r++
			if r > 6 {
				break
			}
		}
	}
	switch {
	case r == 0:
		return "Done"
	case r <= 2:
		return "Low"
	case r <= 4:
		return "Medium"
	case r <= 6:
		return "High"
	default:
		return "Critical"
	}
}
