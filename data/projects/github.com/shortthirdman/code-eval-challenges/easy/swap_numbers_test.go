package main

import (
	"strings"
	"testing"
)

func TestSwapNumbers(t *testing.T) {
	for k, v := range map[string]string{
		"4Always0 5look8 4on9 7the2 4bright8 9side7 3of8 5life5": "0Always4 8look5 9on4 2the7 8bright4 7side9 8of3 5life5",
		"5Nobody5 7expects3 5the4 6Spanish4 9inquisition0":       "5Nobody5 3expects7 4the5 4Spanish6 0inquisition9"} {
		if r := swapNumbers(k); r != v {
			t.Errorf("failed: swapNumbers %s is %s, got %s",
				k, v, r)
		}
	}
}

func swapNumbers(q string) string {
	s := strings.Fields(q)
	for ix, i := range s {
		s[ix] = string(i[len(i)-1]) + i[1:len(i)-1] + string(i[0])
	}
	return strings.Join(s, " ")
}
