package act_test

import (
	"testing"

	act "github.com/cristaloleg/go-act"
)

func TestStack(t *testing.T) {
	s := act.NewStack()
	s.Push() <- 10
	s.Push() <- 20

	v := <-s.Pop()
	if v != 20 {
		t.Errorf("expected 20, got %v", v)
	}

	v = <-s.Pop()
	if v != 10 {
		t.Errorf("expected 10, got %v", v)
	}
}
