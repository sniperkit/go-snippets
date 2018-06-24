package main

import (
	"regexp"
	"strings"
	"testing"
)

type tuple struct {
	p, q string
}

func TestSeqTrans(t *testing.T) {
	for k, v := range map[tuple]bool{
		tuple{"1010", "AAAAABBBBAAAA"}:             true,
		tuple{"00", "AAAAAA"}:                      true,
		tuple{"01001110", "AAAABAAABBBBBBAAAAAAA"}: true,
		tuple{"1100110", "BBAABABBA"}:              false} {
		if r := seqTrans(k.p, k.q); r != v {
			t.Errorf("failed: seqTrans %s %s is %t, got %t",
				k.p, k.q, v, r)
		}
	}
}

func seqTrans(p, q string) bool {
	p = "^" + strings.Replace(strings.Replace(p, "0", "A+", -1), "1", "((A+)|(B+))", -1) + "$"
	seqRegex := regexp.MustCompile(p)
	return seqRegex.MatchString(q)
}
