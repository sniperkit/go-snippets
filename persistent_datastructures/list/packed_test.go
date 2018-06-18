package list

import "testing"

func TestPacked1(t *testing.T) {
	var p List
	p = packed{1, 2, 3, 4, 5}
	for i := 1; i <= 5; i++ {
		v := p.Value()
		p = p.Next()
		if v != i {
			t.Fatalf("%dth value not equal %d, got %d", i, i, v)
		}
	}
}
