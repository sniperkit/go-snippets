package list

import "testing"

func TestUglyList1(t *testing.T) {
	n := Node{1, nil}
	if m, ok := Value(n); !ok || m != n.Value {
		t.Fatalf("node value should equal 1, but we got %d", m)
	}
}

func TestUglyList2(t *testing.T) {
	n := End{}
	if m, ok := Value(n); ok || m != 0 {
		t.Fatalf("end value should equal 0", m)
	}
}

func testUglyList3(t *testing.T) {
	n := struct{}{}
	Value(n)
}
