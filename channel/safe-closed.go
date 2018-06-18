package main

import (
	"fmt"
	"sync"
)

type T int

// method 1
func SafeClosed(c chan T) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	// assume c != nil here
	close(c) // panic if c is closed
	return true
}

// method 2
type Foo struct {
	C    chan T
	once sync.Once
}

func NewFoo() *Foo {
	return &Foo{C: make(chan T)}
}

func (f *Foo) SafeClosed() {
	f.once.Do(func() {
		close(f.C)
	})
}

func main() {
	ch := make(chan T)
	fmt.Println("safe close channel: ", SafeClosed(ch))
	fmt.Println("safe close channel again: ", SafeClosed(ch))

	f := NewFoo()
	f.SafeClosed()
	f.SafeClosed()
}
