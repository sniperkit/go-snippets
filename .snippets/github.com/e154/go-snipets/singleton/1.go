package main

import "fmt"

type single struct {
	O interface{};
}

var instantiated *single = nil

func New() *single {
	if instantiated == nil {
		instantiated = new(single);
	}
	return instantiated;
}

func main() {
	A := New()
	B := New()
	C := New()

	if A == B && A == C {
		fmt.Println("A == B && A == C")
	}
}

