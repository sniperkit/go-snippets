package main

import "fmt"

type Iface interface {
	DoStuff()
}

type Impl struct{}

func (i *Impl) DoStuff() {
	fmt.Println("Stuff")
}

type SomeType struct{
	Iface
}

func (t *SomeType) Do() {
	t.DoStuff()
}

func main() {
	t := SomeType{&Impl{}}
	t.Do()
}
