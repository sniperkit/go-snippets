package binarytrees

import (
	"github.com/sahilm/go-tour/stack"
	"golang.org/x/tour/tree"
)

// Walk Iteratively traverses binary tree t in-order
// and publishes found nodes on channel c.
func Walk(t *tree.Tree, c chan int) {
	defer close(c)

	s := stack.New()
	for !s.IsEmpty() || t != nil {
		if t != nil {
			s.Push(t)
			t = t.Left
		} else {
			t = pop(&s)
			c <- t.Value
			t = t.Right
		}
	}
}

// Same Is true if the traversal of t1 and t2
// yield the same nodes else is false.
func Same(t1, t2 *tree.Tree) bool {
	c1 := make(chan int)
	go Walk(t1, c1)
	c2 := make(chan int)
	go Walk(t2, c2)

	for {
		v1, ok1 := <-c1
		v2, ok2 := <-c2

		if v1 != v2 || ok1 && !ok2 || !ok1 && ok2 {
			return false
		}

		if !ok1 && !ok2 {
			break
		}
	}

	return true
}

func pop(s *stack.Stack) *tree.Tree {
	v, err := s.Pop()

	if err != nil {
		panic(err)
	}

	return v.(*tree.Tree)
}
