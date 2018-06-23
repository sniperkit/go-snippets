package stack

import "testing"

func TestNewStack(t *testing.T) {
	s := New()
	if s == nil {
		t.Error("Stack must not be nil")
	}
}

func TestPush(t *testing.T) {
	s := New()

	s.Push(1)
	if s.Top().(int) != 1 {
		t.Error("The top element must be 1")
	}
	s.Push(2)
	if s.Top().(int) != 2 {
		t.Error("The top element must be 2")
	}
}

func TestPop(t *testing.T) {
	s := New()

	s.Push(1)
	s.Push(2)

	v := s.Pop()
	if v.(int) != 2 {
		t.Error("The value must be 2")
	}
	v = s.Pop()
	if v.(int) != 1 {
		t.Error("The value must be 1")
	}

	if s.Top() != nil {
		t.Error("The top element must be nil")
	}
}
