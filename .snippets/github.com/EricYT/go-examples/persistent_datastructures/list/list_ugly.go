package list

import "fmt"

// inspired by https://blog.merovius.de/2018/02/25/persistent_datastructures_with_go.html

type List1 interface {
	// We use an unexported marker-method. As nothing outside the current package
	// can implement this unexported method, we get control over all
	// implementions of List and can thus de-facto close the set of possible
	// types.
	list()
}

type Node struct {
	Value int
	Next  List1
}

func (Node) list() {}

type End struct{}

func (End) list() {}

func Value(l interface{}) (v int, ok bool) {
	switch l := l.(type) {
	case Node:
		return l.Value, true
	case End:
		return 0, false
	default:
		// This should never happen. Someone violated our sum-type assumption.
		panic(fmt.Errorf("unknow type %T", l))
	}
}
