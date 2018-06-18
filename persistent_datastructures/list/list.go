package list

type List interface {
	// Value returns the current value of the List
	Value() int
	// Next returns the tail of the List, or nil, if this is the last node.
	Next() List
}

type node struct {
	value int
	next  List
}

func (n node) Value() int {
	return n.value
}

func (n node) Next() List {
	return n.next
}

func NewNode(v int) List {
	return node{v, nil}
}

func Prepend(l List, v int) List {
	return node{v, l}
}
