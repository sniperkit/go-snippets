package main

import (
	"fmt"
	"testing"
)

func TestTrickOrTreat(t *testing.T) {
	for k, v := range map[string]uint{
		"Vampires: 1, Zombies: 1, Witches: 1, Houses: 1":         4,
		"Vampires: 3, Zombies: 2, Witches: 1, Houses: 10":        36,
		"Vampires: 100, Zombies: 100, Witches: 100, Houses: 100": 400,
		"Vampires: 100, Zombies: 100, Witches: 100, Houses: 0":   0,
		"Vampires: 0, Zombies: 0, Witches: 0, Houses: 100":       0} {
		var a, b, c, d uint
		fmt.Sscanf(k,
			"Vampires: %d, Zombies: %d, Witches: %d, Houses: %d",
			&a, &b, &c, &d)
		if r := trickOrTreat(a, b, c, d); r != v {
			t.Errorf("failed: trickOrTreat %s is %d, got %d",
				k, v, r)
		}
	}
}

func BenchmarkTrickOrTreat(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var a, b, c, d uint
		j := uint(i)
		a, j = j%101, j/101
		b, j = j%101, j/101
		c, j = j%101, j/101
		d, j = j%101, j/101
		trickOrTreat(a, b, c, d)
	}
}

func trickOrTreat(v, z, w, h uint) uint {
	if v+z+w == 0 {
		return 0
	}
	return (v*3 + z*4 + w*5) * h / (v + z + w)
}
