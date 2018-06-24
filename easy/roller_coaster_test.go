package main

import "testing"

func TestRollerCoaster(t *testing.T) {
	for k, v := range map[string]string{
		"To be, or not to be: that is the question.":   "To Be, Or NoT tO bE: tHaT iS tHe QuEsTiOn.",
		"Whether 'tis nobler in the mind to suffer":    "WhEtHeR 'tIs NoBlEr In ThE mInD tO sUfFeR",
		"The slings and arrows of outrageous fortune,": "ThE sLiNgS aNd ArRoWs Of OuTrAgEoUs FoRtUnE,",
		"Or to take arms against a sea of troubles,":   "Or To TaKe ArMs AgAiNsT a SeA oF tRoUbLeS,",
		"And by opposing end them? To die: to sleep.":  "AnD bY oPpOsInG eNd ThEm? To DiE: tO sLeEp."} {
		if r := rollerCoaster(k); r != v {
			t.Errorf("failed: rollerCoaster %s is %s, got %s",
				k, v, r)
		}
	}
}

func rollerCoaster(q string) string {
	var u bool
	r := make([]byte, len(q))
	for ix, i := range q {
		if (i >= 'a' && i <= 'z') || (i >= 'A' && i <= 'Z') {
			u = !u
			if u {
				r[ix] = byte(i) & 223
			} else {
				r[ix] = byte(i) | 32
			}
		} else {
			r[ix] = byte(i)
		}
	}
	return string(r)
}
