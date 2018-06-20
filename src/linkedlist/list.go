package linkedlist

import (
	"fmt"
)

type List struct {
	root *chain
}

func NewList() *List {
	return &List{}
}

func (l *List) Insert(key, value interface{}) {
	if l.root == nil {
		l.root = &chain{key: key, value: value}
	} else {
		l.root.insert(key, value)
	}
}

func (l *List) Delete(value interface{}) bool {
	if l.root == nil {
		return false
	} else {
		ok := false
		l.root, ok = l.root.delete(value)
		return ok
	}
}

func (l *List) Find(key interface{}) interface{} {
	if l.root == nil {
		return nil
	} else {
		return l.root.find(key)
	}
}

func (l *List) Length() int {
	if l.root == nil {
		return 0
	} else {
		return l.root.length()
	}
}

func (l *List) Contains(value interface{}) bool {
	if l.root == nil {
		return false
	} else {
		return l.root.contains(value)
	}
}

func (l *List) ForEach(f func(interface{}, interface{})) {
	if l.root != nil {
		l.root.iterate(f)
	}
}

func (l *List) String() string {
	if l.root == nil {
		return "[]"
	} else {
		return "[" + l.root.String() + "]"
	}
}

type chain struct {
	key, value interface{}
	next       *chain
}

func (c *chain) insert(key, value interface{}) {
	node := c
	for node != nil {
		if node.next == nil {
			node.next = &chain{key: key, value: value}
			return
		}
		node = node.next
	}
}

func (c *chain) delete(key interface{}) (*chain, bool) {
	node := c
	last := node
	for node != nil {
		if node.key == key {
			if node == c {
				return c.next, true
			} else {
				last.next = node.next
				return c, true
			}
		}
		last = node
		node = node.next
	}
	return c, false
}

func (c *chain) find(key interface{}) interface{} {
	node := c
	for node != nil {
		if node.value == key {
			return node.value
		}
		node = node.next
	}
	return nil
}

func (c *chain) length() int {
	node := c
	length := 0
	for node != nil {
		length += 1
		node = node.next
	}
	return length
}

func (c *chain) contains(key interface{}) bool {
	chain := c
	for chain != nil {
		if chain.value == key {
			return true
		}
		chain = chain.next
	}
	return false
}

func (c *chain) String() string {
	chain := c
	result := ""
	for chain != nil {
		result += fmt.Sprintf("%v", chain.key) + "=" + fmt.Sprintf("%v", chain.value)
		if chain.next != nil {
			result += ","
		}
		chain = chain.next
	}
	return result
}

func (c *chain) iterate(f func(interface{}, interface{})) {
	chain := c
	for chain != nil {
		f(chain.key, chain.value)
		chain = chain.next
	}
}
