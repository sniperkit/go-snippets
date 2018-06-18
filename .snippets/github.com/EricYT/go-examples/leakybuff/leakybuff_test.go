package leakybuff

import (
	"testing"
)

func TestLeakyBuffSize(t *testing.T) {
	bufsize := 5
	cap := 2
	lb := NewLeakyBuff(bufsize, cap)

	// buf size
	buf := lb.Get()
	if len(buf) != bufsize {
		t.Fatalf("leaky buffer unit size not equal 2 but %d", len(buf))
	}
}

func TestLeakyBuffPut(t *testing.T) {
	lb := NewLeakyBuff(5, 2)

	buf := lb.Get()
	if buf == nil {
		t.Fatalf("leaky buffer should get a fresh unit")
	}

	buf = []byte("hello")

	lb.Put(buf)

	bufLater := lb.Get()
	if string(bufLater) != "hello" {
		t.Fatalf("should got same buffer but got data(%s)", string(bufLater))
	}
}

func TestLeakyBuffGet(t *testing.T) {
	lb := NewLeakyBuff(5, 2)

	buf := lb.Get()
	lb.Put(buf)

	bufNew := lb.Get()

	bufPtr := &buf[0]
	bufNewPtr := &bufNew[0]
	if bufPtr != bufNewPtr {
		t.Fatalf("should got same buffer")
	}

	bufNew1 := lb.Get()
	bufNewPtr1 := &bufNew1[0]
	if bufPtr == bufNewPtr1 {
		t.Fatalf("should got a new buffer")
	}

}
