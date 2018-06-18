package math

import "testing"

// Use unsigned int of four bits
// 0 0 0 0

// 0 0 1 0
// 0 1 0 0
// 0 1 1 0
// 1 0 0 0
// 1 0 1 0
// 1 1 0 0
// 1 1 1 0
// 0 0 0 0  (bingo)

func TestOverflowUint8Even(t *testing.T) {
	var i uint8 = 0

	for {
		before := i
		i += 2
		if i == i%2 {
			t.Errorf("A number start from zero has a increase step two. overflow when it is: %d before: %d", i, before)
			return
		}
	}
}

func TestOverflowUint8Odd(t *testing.T) {
	var i uint8 = 1

	for {
		before := i
		i += 2
		if i == i%2 {
			t.Errorf("A number start from one has a increase step two. overflow when it is: %d before: %d", i, before)
			return
		}
	}
}
