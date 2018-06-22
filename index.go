package skiplist

import (
	"sync/atomic"
	"unsafe"
)

type index struct {
	node  *node
	right *index
	down  *index
	level int
}

//create header
func newHeader(key []byte, value unsafe.Pointer, next *node, down, right *index, level int) *index {
	return &index{
		node: &node{
			key:   key,
			value: value,
			next:  next,
		},
		down:  down,
		right: right,
		level: level,
	}
}

//create index
func newIndex(n *node, down, right *index) *index {
	return &index{
		node:  n,
		down:  down,
		right: right,
	}
}

//delete marked node
func (i *index) deleteMarkedNode(succ *index) bool {
	return atomic.LoadInt32(&i.node.marked) == 0 && i.casRight(succ, succ.right)
}

//cas right index
func (i *index) casRight(cmp, right *index) bool {
	return atomic.CompareAndSwapPointer((*unsafe.Pointer)(unsafe.Pointer(&i.right)), unsafe.Pointer(cmp), unsafe.Pointer(right))
}

//add index
func (i *index) addIndex(succ, newSucc *index) bool {
	newSucc.right = succ

	return atomic.LoadInt32(&i.node.marked) == 0 && i.casRight(succ, newSucc)
}
