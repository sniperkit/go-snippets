package twothree

import (
	"testing"
	"math/rand"
)

func randTree(n int) ([]int, *tree) {
	tree := NewTree()
	values := make([]int, 200)
	for i := 0; i < 200; i++ {
		values[i] = rand.Intn(500) - 250
	}

	for i := range values {
		tree.Insert(values[i])
	}
	return values, tree
}

func TestTree_Length(t *testing.T) {
	values,tree := randTree(200)
	RemoveDuplicates(&values)
	tLen, vLen := tree.Length(), len(values)
	if tLen != vLen {
		t.Error(
			"For", tree,
			"expected", vLen,
			"got", tLen,
		)
	}
}

func TestTree_Insert(t *testing.T) {
	values,tree := randTree(200)
	for i := range values {
		if !tree.Member(values[i]) {
			t.Error(
				"For", values[i],
				"expected", true,
				"got", false,
			)
		}
	}
}

func RemoveDuplicates(xs *[]int) {
	found := make(map[int]bool)
	j := 0
	for i, x := range *xs {
		if !found[x] {
			found[x] = true
			(*xs)[j] = (*xs)[i]
			j++
		}
	}
	*xs = (*xs)[:j]
}
