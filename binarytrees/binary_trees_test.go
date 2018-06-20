package binarytrees_test

import (
	"testing"

	"github.com/sahilm/go-tour/binarytrees"
	"golang.org/x/tour/tree"
)

func TestWalk(t *testing.T) {
	cases := []struct {
		in   *tree.Tree
		want []int
	}{
		{tree.New(1), []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}},
		{tree.New(2), []int{2, 4, 6, 8, 10, 12, 14, 16, 18, 20}},
	}

	for _, c := range cases {
		ch := make(chan int)
		go binarytrees.Walk(c.in, ch)
		for i := 0; i < 10; i++ {
			got := <-ch
			if got != c.want[i] {
				t.Errorf("got: %v, want: %v", got, c.want[i])
			}
		}
	}
}

func TestSame(t *testing.T) {
	cases := []struct {
		t1, t2 *tree.Tree
		want   bool
	}{
		{tree.New(1), tree.New(1), true},
		{tree.New(2), tree.New(1), false},
		{tree.New(5), tree.New(5), true},
	}

	for _, c := range cases {
		got := binarytrees.Same(c.t1, c.t2)
		if got != c.want {
			t.Errorf("got: %v, want: %v", got, c.want)
		}
	}
}
