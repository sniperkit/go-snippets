package main

import "testing"

func TestNrc(t *testing.T) {
	for k, v := range map[string]string{
		"yellow":     "y",
		"tooth":      "h",
		"blablaG":    "G",
		"Gagaga":     "G",
		"1234512345": "",
		"in nova fert animus mutatas dicere formas corpora di coeptis nam vos mutastis et illas adspirate meis primaque ab origine mundi": "q"} {
		if r := nrc(k); r != v {
			t.Errorf("failed: nrc %s is %s, got %s",
				k, v, r)
		}
	}
}

func nrc(q string) string {
	s := make(map[rune]int)
	for _, i := range q {
		s[i]++
	}
	for _, i := range q {
		if s[i] == 1 {
			return string(i)
		}
	}
	return ""
}
