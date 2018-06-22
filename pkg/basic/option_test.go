package basic

import (
	"fmt"
	"testing"
)

type Foo struct {
	A         int
	verbosity int
}

type option func(f *Foo) option

func Verbosity(v int) option {
	return func(f *Foo) option {
		previous := f.verbosity
		f.verbosity = v
		return Verbosity(previous)
	}
}

func (f *Foo) Option(opts ...option) (previous option) {
	for _, opt := range opts {
		previous = opt(f)
	}
	return previous
}

func DoSomethingVerbosely(foo *Foo, verbosity int) {
	// Could combine the next two lines,
	// with some loss of readability.
	prev := foo.Option(Verbosity(verbosity))
	defer foo.Option(prev)
	// ... do some stuff with foo under high verbosity.
}

func TestOption(t *testing.T) {
	var opt option
	fmt.Println(opt)
}
