package maxmin

import (
	"math/rand"
	"testing"
	"time"
)

func TestSamples(t *testing.T) {

	is := []int{}
	max, min := IndexOfMaxMin(IntSlice(is))
	if (max != -1) || (min != -1) {
		t.Fatal("max and min must = -1")
	}

	is = []int{-7}
	max, min = IndexOfMaxMin(IntSlice(is))
	if (max != 0) || (min != 0) {
		t.Fatal("max and min must = 0")
	}

	is = []int{1, 2, 3, 4, 5}
	max, min = IndexOfMaxMin(IntSlice(is))
	if (max != 4) || (min != 0) {
		t.Fatal("wrong max or min")
	}

	is = []int{9, 8, 7, 6, 5, 4}
	max, min = IndexOfMaxMin(IntSlice(is))
	if (max != 0) || (min != 5) {
		t.Fatal("wrong max or min")
	}

	is = []int{0, 1, -5, -5, -5, -5, 3, 3, 3}
	max, min = IndexOfMaxMin(IntSlice(is))
	if (max != 6) || (min != 2) {
		t.Fatal("wrong max or min")
	}
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func randRange(r *rand.Rand, min, max int) int {
	if min > max {
		min, max = max, min
	}
	return min + r.Intn(max-min)
}

func TestRand(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for j := 0; j < 1000; j++ {
		is := make([]int, 1+r.Intn(100))
		for i := range is {
			is[i] = randRange(r, -200, +201)
		}
		max, min := is[0], is[0]
		for i := range is {
			max = maxInt(max, is[i])
			min = minInt(min, is[i])
		}
		imax, imin := IndexOfMaxMin(IntSlice(is))
		if is[imax] != max {
			t.Fatal("max wrong")
		}
		if is[imin] != min {
			t.Fatal("min wrong")
		}
	}
}
