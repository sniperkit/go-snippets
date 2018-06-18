package main

type SomeInterface interface {
	someFunction(someArg int)
}

type someStructA struct {
	someState int
}

type someStructB struct {
	someState int
}

func (s someStructA) someFunction(someArg int) {
	s.someState = someArg
}

func (s *someStructB) someFunction(someArg int) {
	s.someState = someArg
}

func main() {
	var a SomeInterface = someStructA{}
	var b SomeInterface = &someStructB{}

	print(a)
	print(b)

	// Чем отличаются (по поведению) следующие два вызова?
	a.someFunction(1)
	b.someFunction(1)
}