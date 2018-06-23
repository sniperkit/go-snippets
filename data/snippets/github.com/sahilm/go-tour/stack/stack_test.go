package stack_test

import (
	"errors"
	"testing"

	"github.com/sahilm/go-tour/stack"
)

func TestStack(t *testing.T) {
	s := stack.New()
	if !s.IsEmpty() {
		t.Errorf("got: %v, want empty stack", s)
	}

	s.Push(1, 2, 3)
	if s.IsEmpty() {
		t.Error("got empty stack, want non-empty stack")
	}

	for _, v := range []int{3, 2, 1} {
		got, err := s.Pop()
		if err != nil {
			t.Errorf("got error: %v, want no error", err)
		}
		if got != v {
			t.Errorf("got: %v, want: %v", got, v)
		}
	}

	_, err := s.Pop()
	want := errors.New("cannot pop an empty stack")
	if err.Error() != want.Error() {
		t.Errorf("got error: %v, want error: %v", err, want)
	}
}
