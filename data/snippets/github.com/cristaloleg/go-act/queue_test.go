package act_test

import (
	"testing"

	act "github.com/cristaloleg/go-act"
)

func TestQueue(t *testing.T) {
	q := act.NewQueue()

	q.Push() <- 10
	q.Push() <- 20

	v := <-q.Pop()
	if v != 10 {
		t.Errorf("expected 10, got %v", v)
	}

	v = <-q.Pop()
	if v != 20 {
		t.Errorf("expected 20, got %v", v)
	}
}
