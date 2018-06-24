package main

import "testing"

type tuple struct {
	q string
	c int
	l bool
}

func TestSlang(t *testing.T) {
	for k, v := range map[tuple]tuple{
		tuple{"a", 0, false}: tuple{"a", 0, false},
		tuple{"a", 0, true}:  tuple{"a", 0, true},
		tuple{".", 0, false}: tuple{".", 0, true},
		tuple{".", 0, true}:  tuple{", yeah!", 1, false},
		tuple{"!", 2, true}:  tuple{", can U believe this?", 3, false},
		tuple{"?", 7, true}:  tuple{". Awesome!", 0, false},
	} {
		if rq, rc, rl := slang(k.q, k.c, k.l); rq != v.q || rc != v.c || rl != v.l {
			t.Errorf("failed: slang %s %d %t is %s %d %t, got %s %d %t",
				k.q, k.c, k.l, v.q, v.c, v.l, rq, rc, rl)
		}
	}
}

var (
	phrases = []string{", yeah!", ", this is crazy, I tell ya.",
		", can U believe this?", ", eh?", ", aw yea.",
		", yo.", "? No way!", ". Awesome!"}
)

func slang(q string, c int, l bool) (string, int, bool) {
	if q == "." || q == "!" || q == "?" {
		if l {
			return phrases[c], (c + 1) % len(phrases), false
		}
		return q, c, true
	}
	return q, c, l
}
