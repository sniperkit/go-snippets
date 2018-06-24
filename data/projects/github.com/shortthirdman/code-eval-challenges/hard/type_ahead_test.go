package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"
)

type tuple struct {
	n int
	q string
}

func TestTypeAhead(t *testing.T) {
	for k, v := range map[tuple]string{
		tuple{2, "the"}:         "lamb,0.375;teacher,0.250;children,0.125;eager,0.125;rule,0.125",
		tuple{3, "the teacher"}: "did,0.500;turned,0.500",
		tuple{4, "And"}:         "everywhere that Mary,0.333;so the teacher,0.333;waited patiently about,0.333",
		tuple{2, "Mary"}:        "did,0.200;had,0.200;loves,0.200;so,0.200;went,0.200",
		tuple{4, "the lamb"}:    "love Mary,0.333;was sure,0.333;you know,0.333"} {
		if r := typeAhead(k.n, k.q); r != v {
			t.Errorf("failed: typeAhead %d %s is %s, got %s",
				k.n, k.q, v, r)
		}
	}
}

const text = `
Mary had a little lamb its fleece was white as snow;
And everywhere that Mary went, the lamb was sure to go. 
It followed her to school one day, which was against the rule;
It made the children laugh and play, to see a lamb at school.
And so the teacher turned it out, but still it lingered near,
And waited patiently about till Mary did appear.
"Why does the lamb love Mary so?" the eager children cry;
"Why, Mary loves the lamb, you know" the teacher did reply."
`

var (
	freq  map[string]int
	words []string
)

func init() {
	words = strings.Fields(text)
	for ix, i := range words {
		words[ix] = strings.Trim(i, ";,.?\"")
	}
}

type ngrams [][]string

func (s ngrams) Len() int { return len(s) }
func (s ngrams) Less(i, j int) bool {
	a, b := strings.Join(s[i], ","), strings.Join(s[j], ",")
	if freq[a] == freq[b] {
		return a < b
	}
	return freq[a] > freq[b]
}
func (s ngrams) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func sliceq(a, b []string) bool {
	for ix, i := range a {
		if i != b[ix] {
			return false
		}
	}
	return true
}

func typeAhead(n int, q string) string {
	var (
		tnl, r [][]string
		p      []string
	)
	x := strings.Fields(strings.TrimSpace(q))

	for i := 0; i <= len(words)-n; i++ {
		temp := make([]string, n)
		for j := 0; j < n; j++ {
			temp[j] = words[i+j]
		}
		tnl = append(tnl, temp)
	}

	for _, i := range tnl {
		f := true
		for j := range x {
			if i[j] != x[j] {
				f = false
				break
			}
		}
		if f {
			r = append(r, i)
		}
	}

	freq = make(map[string]int)
	for _, i := range r {
		rs := strings.Join(i, ",")
		if _, f := freq[rs]; f {
			freq[rs]++
		} else {
			freq[rs] = 1
		}
	}

	sort.Sort(ngrams(r))
	for ix, i := range r {
		if ix == 0 || !sliceq(i, r[ix-1]) {
			p = append(p, strings.Join(i[len(x):], " ")+fmt.Sprintf(",%.3f", float32(freq[strings.Join(i, ",")])/float32(len(r))))
		}
	}
	return strings.Join(p, ";")
}
