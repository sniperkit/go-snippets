package atomic_test

import (
	"testing"

	. "github.com/cristaloleg/funf/atomic"
)

func TestAtomic(t *testing.T) {
	a := New(fn)
	if a == nil {
		t.Error("cannot instantiate")
	}

	if value := a.Get(); value != nil {
		t.Errorf("expected default value nil, got %v", value)
	}

	a.Update("test")

	if value := a.Get(); value != "test" {
		t.Errorf("expected %v, got %v", "test", value)
	}
}

func TestAtomic_Min(t *testing.T) {
	a := NewMin()
	if a == nil {
		t.Error("cannot instantiate")
	}

	if value := a.Get(); value != nil {
		t.Errorf("expected default value nil, got %v", value)
	}

	a.Update(10)
	if value := a.Get(); value != 10 {
		t.Errorf("expected %v, got %v", 10, value)
	}

	a.Update(12)
	if value := a.Get(); value != 10 {
		t.Errorf("expected %v, got %v", 10, value)
	}

	a.Update(7)
	if value := a.Get(); value != 7 {
		t.Errorf("expected %v, got %v", 7, value)
	}
}

func TestAtomic_Max(t *testing.T) {
	a := NewMax()
	if a == nil {
		t.Error("cannot instantiate")
	}

	if value := a.Get(); value != nil {
		t.Errorf("expected default value nil, got %v", value)
	}

	a.Update(10)
	if value := a.Get(); value != 10 {
		t.Errorf("expected %v, got %v", 10, value)
	}

	a.Update(7)
	if value := a.Get(); value != 10 {
		t.Errorf("expected %v, got %v", 10, value)
	}

	a.Update(12)
	if value := a.Get(); value != 12 {
		t.Errorf("expected %v, got %v", 12, value)
	}
}

func fn(x, y interface{}) bool {
	return true
}
