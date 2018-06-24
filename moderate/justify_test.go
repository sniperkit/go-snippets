package main

import (
	"strings"
	"testing"
)

func TestJustify(t *testing.T) {
	for k, v := range map[string]string{
		"Hello, World!": "Hello, World!",
		"The precise 50-digits value of Pi is 3.14159265358979323846264338327950288419716939937510.": `The         precise         50-digits        value        of        Pi        is
3.14159265358979323846264338327950288419716939937510.`,
		"But he who would be a great man ought to regard, not himself or his interests, but what is just, whether the just act be his own or that of another. Next as to habitations. Such is the tradition.": `But  he  who would be a great man ought to regard, not himself or his interests,
but what is just, whether the just act be his own or that of another. Next as to
habitations. Such is the tradition.`} {
		if r := justify(k); r != v {
			t.Errorf("failed: justify\n%s\n is\n%s\n got\n%s",
				k, v, r)
		}
	}
}

func justi(s string) string {
	if len(s) == 80 {
		return s
	}
	ws := strings.Fields(s)
	l := 80 - len(strings.Join(ws, ""))
	for i := 0; i < l%(len(ws)-1); i++ {
		ws[i] = ws[i] + " "
	}
	return strings.Join(ws, strings.Repeat(" ", l/(len(ws)-1)))
}

func maxwords(s string) int {
	if len(s) > 81 {
		return maxwords(s[:81])
	}
	return strings.LastIndex(s, " ")
}

func justif(s string) []string {
	t := strings.TrimLeft(s, " ")
	if len(t) <= 80 {
		return []string{t}
	}
	m := maxwords(t)
	return append([]string{justi(t[:m])}, justif(t[m:])...)
}

func justify(s string) string {
	if len(s) <= 80 {
		return s
	}
	return strings.Join(justif(s), "\n")
}
