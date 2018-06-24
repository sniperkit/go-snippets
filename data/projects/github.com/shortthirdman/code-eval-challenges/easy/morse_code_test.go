package main

import "testing"

func TestMorseCode(t *testing.T) {
	for k, v := range map[string]string{
		".- ...- ..--- .-- .... .. . -.-. -..-  ....- .....": "AV2WHIECX 45",
		"-... .... ...--":                                    "BH3",
		"-----":                                              "0"} {
		if r := morseCode(k); r != v {
			t.Errorf("failed: morseCode %s is %s, got %s",
				k, v, r)
		}
	}
}

const morse = "ETIANMSURWDKGOHVF L PJBXCYZQ  54 3   2       16       7   8 90"

func morseCode(q string) string {
	var (
		m int = 1
		r []byte
	)
	for _, i := range q {
		switch {
		case i == '.':
			m <<= 1
		case i == '-':
			m = (m << 1) + 1
		case m == 1:
			r = append(r, ' ')
		default:
			if m < 64 {
				r = append(r, morse[m-2])
			}
			m = 1
		}
	}
	if m > 1 && m < 64 {
		r = append(r, morse[m-2])
	}
	return string(r)
}
