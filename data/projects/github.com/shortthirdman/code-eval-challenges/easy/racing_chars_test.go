package main

import (
	"strings"
	"testing"
)

type tuple struct {
	n int
	q string
}

func TestRacing(t *testing.T) {
	for k, v := range map[tuple]tuple{
		tuple{9, "########C_##"}:  tuple{8, "########/_##"},
		tuple{8, "#######_####"}:  tuple{7, "#######/####"},
		tuple{7, "######_#C###"}:  tuple{8, "######_#\\###"},
		tuple{8, "#######_C###"}:  tuple{8, "#######_|###"},
		tuple{10, "###########_"}: tuple{11, "###########\\"},
		tuple{11, "###########_"}: tuple{11, "###########|"},
		tuple{1, "_###########"}:  tuple{0, "/###########"},
		tuple{0, "_###########"}:  tuple{0, "|###########"}} {
		if rn, rq := racing(k.n, k.q); rn != v.n || rq != v.q {
			t.Errorf("failed: racing %d %s is %d %s, got %d %s",
				k.n, k.q, v.n, v.q, rn, rq)
		}
	}
}

func racing(n int, q string) (int, string) {
	p := n
	a, b := p-1, p+2
	if a < 0 {
		a = 0
	}
	if b > len(q) {
		b = len(q)
	}
	c := strings.Index(q[a:b], "C")
	if c == -1 {
		c = a + strings.Index(q[a:b], "_")
	} else {
		c += a
	}
	if c < p {
		q = q[:c] + "/" + q[c+1:]
	} else if c == p {
		q = q[:c] + "|" + q[c+1:]
	} else {
		q = q[:c] + `\` + q[c+1:]
	}
	return c, q
}
