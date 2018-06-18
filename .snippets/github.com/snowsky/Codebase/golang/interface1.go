package main

import "fmt"

type Dog struct {
	name string
}
type BDog struct {
	Dog
	name string
}
type Funny interface {
	callMyName()
	getName() string
}

// func (this *BDog) callMyName() {
// 	fmt.Printf("my name is %q\n", this.getName())
// }
// func (this *BDog) getName() string {
// 	return this.name
// }

func (this *Dog) callMyName() {
	fmt.Printf("my name is %q\n", this.getName())
}
func (this *Dog) getName() string {
	return this.name
}
func main() {
	b := new(BDog)
	b.name = "this is a BDog name"
	b.Dog.name = "this is a Dog name"
	b.callMyName()
}
