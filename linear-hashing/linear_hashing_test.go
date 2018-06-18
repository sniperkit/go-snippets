package linear_hashing

import (
	"fmt"
	"testing"
)

func TestLinearHashing(t *testing.T) {
	h := NewLinearHashing(4, 3)

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key#%d", i)
		value := fmt.Sprintf("value#%d", i)
		h.Set(key, value)
	}

	for i := 0; i < 20; i++ {
		key := fmt.Sprintf("key#%d", i)
		value := fmt.Sprintf("value#%d", i)
		val, ok := h.Get(key)
		if !ok || val != value {
			t.Fatalf("get key: %s value: %s not match %s", key, val, value)
		}
	}
}

func BenchmarkSet(b *testing.B) {
	h := NewLinearHashing(4, 5)
	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key#%d", i)
		value := fmt.Sprintf("value#%d", i)
		h.Set(key, value)
	}
}

func BenchmarkGet(b *testing.B) {
	h := NewLinearHashing(4, 5)
	for i := 0; i < 100000; i++ {
		key := fmt.Sprintf("key#%d", i)
		value := fmt.Sprintf("value#%d", i)
		h.Set(key, value)
	}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		key := fmt.Sprintf("key#%d", i)
		h.Get(key)
	}
}
