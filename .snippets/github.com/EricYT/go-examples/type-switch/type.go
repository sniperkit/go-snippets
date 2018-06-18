package main

import "fmt"

type Gopher interface {
	Go() string
}

type Foo int

func (f Foo) Go() string { return "foo" }

type Bar int

func (b Bar) Go() string { return "bar" }

func FooCome(g Gopher) (*Foo, bool) {
	switch f := g.(type) {
	case *Foo:
		return f, true
	default:
		return nil, false
	}
}

func BarCome(g Gopher) (*Bar, bool) {
	switch b := g.(type) {
	case *Bar:
		return b, true
	default:
		return nil, false
	}
}

func main() {
	fmt.Println("type-switch-go")

	var f *Foo = new(Foo)
	var b *Bar = new(Bar)

	fmt.Println("foo come foo")
	if foo, ok := FooCome(f); ok {
		fmt.Println(foo.Go())
	}

	fmt.Println("foo come bar")
	if bar, ok := BarCome(f); ok {
		fmt.Println(bar.Go())
	}

	fmt.Println("bar come foo")
	if foo, ok := FooCome(b); ok {
		fmt.Println(foo.Go())
	}

	fmt.Println("bar come bar")
	if bar, ok := BarCome(b); ok {
		fmt.Println(bar.Go())
	}
}
