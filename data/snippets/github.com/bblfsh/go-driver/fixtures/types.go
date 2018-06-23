package fixtures

import "bytes"

type A int

type B = A

type Buf bytes.Buffer

type Arr []string

type Obj struct {
	io.Writer
	_      int `json:"-" foo:"int"`
	A, B   string
	Field  *[]Obj
	Inline struct {
		Name string
	}
}

type Void interface{}

type Node interface {
	Void
	IsNode()
}

type LongMap map[string]map[int][]string

type Chan chan int
type ChanIn chan<- int
type ChanOut <-chan struct{}
