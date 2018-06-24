package main

import (
	"sort"
	"strings"
	"testing"
)

func TestCodeCombinations(t *testing.T) {
	for k, v := range map[string]uint{
		"**** | *co* | *de* | ****": 1,
		"codx | decx":               2,
		"co | dx":                   0} {
		if r := codeCombinations(k); r != v {
			t.Errorf("failed: codeCombinations %s is %d, got %d",
				k, v, r)
		}
	}
}

func codeCombinations(q string) (r uint) {
	s := strings.Split(q, " | ")
	for i := 1; i < len(s); i++ {
		for j := 1; j < len(s[0]); j++ {
			code := []string{string(s[i][j]), string(s[i][j-1]), string(s[i-1][j]), string(s[i-1][j-1])}
			sort.Strings(code)
			if strings.Join(code, "") == "cdeo" {
				r += 1
			}
		}
	}
	return r
}
