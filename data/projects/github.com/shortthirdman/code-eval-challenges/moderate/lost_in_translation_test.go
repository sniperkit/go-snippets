package main

import (
	"strings"
	"testing"
)

func TestMf(t *testing.T) {
	for k, v := range map[string]string{
		"nug": "bjv",
		"rbc vjnmkf kd yxyqci na rbc zjkfoscdd ew rbc ujllmcp":   "the public is amazed by the quickness of the juggler",
		"tc rbkso rbyr ejp mysljylc kd kxveddknmc re jsicpdrysi": "we think that our language is impossible to understand",
		"de kr kd eoya kw aej icfkici re zjkr":                   "so it is okay if you decided to quit",
		"chfymknjp":                                              "excalibur"} {
		if r := strings.Map(mf, k); r != v {
			t.Errorf("failed: map mf %s is %s, got %s",
				k, v, r)
		}
	}
}

var m map[rune]rune

func init() {
	codel := `nug
rbc vjnmkf kd yxyqci na rbc zjkfoscdd ew rbc ujllmcp
tc rbkso rbyr ejp mysljylc kd kxveddknmc re jsicpdrysi
de kr kd eoya kw aej icfkici re zjkr`
	decol := `bjv
the public is amazed by the quickness of the juggler
we think that our language is impossible to understand
so it is okay if you decided to quit`

	m = map[rune]rune{' ': ' '}
	var misscod, missdec []rune
	for i := 'a'; i <= 'z'; i++ {
		ix := strings.IndexRune(codel, i)
		if ix != -1 {
			m[i] = rune(decol[ix])
		} else {
			misscod = append(misscod, i)
		}

		ix = strings.IndexRune(decol, i)
		if ix == -1 {
			missdec = append(missdec, i)
		}
	}
	for ix, i := range misscod {
		m[i] = missdec[ix]
	}
}

func mf(a rune) rune {
	return m[a]
}
