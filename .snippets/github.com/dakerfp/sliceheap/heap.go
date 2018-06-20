package sliceheap

import (
	"reflect"
)

type SliceHeap struct {
	slice reflect.Value
	less  func(a, b interface{}) bool
}

func (s *SliceHeap) Len() int {
	return s.slice.Len()
}

func (s *SliceHeap) Less(i, j int) bool {
	a := s.slice.Index(i)
	b := s.slice.Index(j)
	return s.less(a.Interface(), b.Interface())
}

func (s *SliceHeap) Swap(i, j int) {
	reflect.Swapper(s.slice.Interface())(i, j)
}

func (s *SliceHeap) Push(x interface{}) {
	s.slice = reflect.Append(s.slice, reflect.ValueOf(x))
}

func (s *SliceHeap) Pop() interface{} {
	n := s.slice.Len()
	old := s.slice.Index(n - 1)
	s.slice = s.slice.Slice(0, n-1)
	return old.Interface()
}

func New(x interface{}, less func(i, j interface{}) bool) *SliceHeap {
	return &SliceHeap{
		slice: reflect.ValueOf(x),
		less:  less,
	}
}
