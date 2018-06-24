package main

import (
	"fmt"
	"testing"
)

var f []float64

func init() {
	f = make([]float64, 5)
}

func TestInCircle(t *testing.T) {
	for k, v := range map[string]bool{
		"Center: (2.12, -3.48); Radius: 17.22; Point: (16.21, -5)":   true,
		"Center: (5.05, -11); Radius: 21.2; Point: (-31, -45)":       false,
		"Center: (-9.86, 1.95); Radius: 47.28; Point: (6.03, -6.42)": true} {
		fmt.Sscanf(k, "Center: (%f, %f); Radius: %f; Point: (%f, %f)",
			&f[0], &f[1], &f[2], &f[3], &f[4])
		if r := inCircle(f[0]-f[3], f[1]-f[4], f[2]); r != v {
			t.Errorf("failed: inCircle %s is %t, got %t",
				k, v, r)
		}
	}
}

func inCircle(a, b, c float64) bool {
	return a*a+b*b <= c*c
}
