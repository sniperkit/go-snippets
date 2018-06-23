package io_queue

import "fmt"

type CircularBuffer struct {
	capacity int
	buffer   []interface{}
	head     int
	tail     int
}

func NewCircularBuffer(cap int) *CircularBuffer {
	cb := new(CircularBuffer)
	cb.capacity = cap
	cb.buffer = make([]interface{}, cb.capacity, cb.capacity)
	cb.head = 0
	cb.tail = 0
	return cb
}

func (cb *CircularBuffer) mask(idx int) int {
	return idx & (cb.capacity - 1)
}

func (cb *CircularBuffer) Empty() bool {
	return cb.head == cb.tail
}

func (cb *CircularBuffer) PushBack(val interface{}) {
	fmt.Printf("before head: %d tail: %d val: %d\n", cb.head, cb.mask(cb.tail), val.(int))
	cb.buffer[cb.mask(cb.tail)] = val
	cb.tail++
	fmt.Printf("after head: %d tail: %d\n", cb.head, cb.tail)
}

func (cb *CircularBuffer) PopFront() interface{} {
	fmt.Printf("+before head: %d tail: %d\n", cb.mask(cb.head), cb.tail)
	if cb.Empty() {
		return nil
	}
	item := cb.buffer[cb.mask(cb.head)]
	cb.head++
	fmt.Printf("+after head: %d tail: %d val: %d\n", cb.head, cb.tail, item.(int))
	return item
}
