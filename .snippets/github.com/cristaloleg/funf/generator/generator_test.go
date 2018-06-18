package generator

import (
	"math/rand"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	fn := func(int64) interface{} { return nil }
	g := New(fn)

	if g == nil {
		t.Error("Should be instantiated")
	}

	size := int64(rand.Int() % 10)
	for i := int64(0); i < size; i++ {
		value := g.Next()
		if value != i {
			t.Error("Should equal to index, was ", value, " instead of ", i)
		}
	}
}

func TestNewBoundedGenerator(t *testing.T) {
	size := int64(10)
	fn := func(x int64) interface{} { return x }
	g := NewBounded(size, fn)

	if g == nil {
		t.Error("Should be instantiated")
	}

	for i := int64(0); i < size; i++ {
		value := g.Next()
		if value != i {
			t.Error("Should equal to index, was ", value, " instead of ", i)
		}
	}
}
