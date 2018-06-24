package main

import (
	"fmt"
	"strings"
	"testing"
)

func TestBuildersTeam(t *testing.T) {
	for k, v := range map[string]uint{
		"1 2 | 6 7 | 7 2 | 1 6 | 2 3":             1,
		"1 2 | 6 7 | 7 2 | 1 6 | 2 3 | 7 8 | 3 8": 2,
		"1 2 | 1 6 | 6 7":                         0} {
		if r := buildersTeam(k); r != v {
			t.Errorf("failed: buildersTeam %s is %d, got %d",
				k, v, r)
		}
	}
}

var (
	a, b     uint
	hor, ver []uint
)

func init() {
	hor = []uint{0x42, 0x84, 0x108, 0x210,
		0x840, 0x1080, 0x2100, 0x4200,
		0x10800, 0x21000, 0x42000, 0x84000,
		0x210000, 0x420000, 0x840000, 0x1080000,
		0x1806, 0x300c, 0x6018,
		0x300c0, 0x60180, 0xc0300,
		0x601800, 0xc03000, 0x1806000,
		0x7000e, 0xe001c,
		0xe001c0, 0x1c00380,
		0x1e0001e}
	ver = []uint{0x6, 0xc, 0x18, 0x30,
		0xc0, 0x180, 0x300, 0x600,
		0x1800, 0x3000, 0x6000, 0xc000,
		0x30000, 0x60000, 0xc0000, 0x180000,
		0x14a, 0x294, 0x528,
		0x2940, 0x5280, 0xa500,
		0x52800, 0xa5000, 0x14a000,
		0x4a52, 0x94a4,
		0x94a40, 0x129480,
		0x118c62}
}

func buildersTeam(q string) (r uint) {
	var h, v uint
	s := strings.Split(q, "|")
	for _, i := range s {
		fmt.Sscanf(i, "%d %d", &a, &b)
		if a > b {
			a, b = b, a
		}
		if a+1 == b {
			h |= 1 << a
		} else {
			v |= 1 << a
		}
	}
	for i := 0; i < 30; i++ {
		if h&hor[i] == hor[i] && v&ver[i] == ver[i] {
			r++
		}
	}
	return r
}
