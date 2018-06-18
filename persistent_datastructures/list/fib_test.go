package list

import "testing"

func TestFib1(t *testing.T) {
	f0 := fib{0, 1}
	if f0.Value() != 0 {
		t.Fatalf("f0 not equal 0, got %d", f0.Value())
	}
	f1 := f0.Next()
	if f1.Value() != 1 {
		t.Fatalf("f1 not equal 1, got %d", f1.Value())
	}
	f2 := f1.Next()
	if f2.Value() != 1 {
		t.Fatalf("f2 not equal 1, got %d", f2.Value())
	}
	f3 := f2.Next()
	if f3.Value() != 2 {
		t.Fatalf("f3 not equal 2, got %d", f3.Value())
	}
	f4 := f3.Next()
	if f4.Value() != 3 {
		t.Fatalf("f4 not equal 3, got %d", f4.Value())
	}
}

func TestFib2(t *testing.T) {
	var f_0 List
	f_0 = fib{0, 1}
	f0, f1 := 0, 1
	for i := 0; i < 20; i++ {
		f0, f1 = f1, f0+f1
		f_0 = f_0.Next()
		if f_0.Value() != f0 {
			t.Fatalf("index %d should equal %d, but got %d", i, f0, f_0.Value())
		}
	}
}
