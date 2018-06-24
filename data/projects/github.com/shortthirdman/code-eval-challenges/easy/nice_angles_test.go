package main

import (
	"fmt"
	"testing"
)

func TestNiceAngles(t *testing.T) {
	for k, v := range map[float64]string{
		330.39991833: "330.23'59\"",
		0.001:        "0.00'03\"",
		14.64530319:  "14.38'43\"",
		0.25:         "0.15'00\"",
		254.16991217: "254.10'11\""} {
		if r := niceAngles(k); r != v {
			t.Errorf("failed: niceAngles %f is %s, got %s",
				k, v, r)
		}
	}
}

func BenchmarkNiceAngles(b *testing.B) {
	for i := 0; i < b.N; i++ {
		niceAngles(float64(i) / 7)
	}
}

func niceAngles(a float64) string {
	return fmt.Sprintf("%d.%02d'%02d\"", int(a),
		int((a-float64(int(a)))*60),
		int((a*60-float64(int(a*60)))*60))
}
