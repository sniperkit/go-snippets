package main

import (
	"fmt"
)

type Oper interface {
	Inc()
	Dec()
}

type Int struct {
	i int
}

func (i *Int) Inc() {
	i.i++
}

func (i *Int) Dec() {
	i.i--
}

type Float struct {
	f float64
}

func (f Float) Inc() {
	f.f += 1.0
}

func (f Float) Dec() {
	f.f -= 1.0
}

func main() {
	//var i Oper = Int{i: 1}
	//var f Oper = Float{f: 3.14}
	var i Oper = &Int{i: 1}
	i.Inc()
	var f Oper = &Float{f: 3.14}
	f.Inc()

	var f1 Oper = Float{f: 3.14}
	f1.Inc()
	fmt.Println(i, f, f1)
}
