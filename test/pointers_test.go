package test

import (
	"fmt"
	"testing"
)

func TestPointer(*testing.T) {
	i, j := 42, 2701

	// point to i
	p := &i

	// read i through the pointer
	fmt.Printf("Read i through the pointer:%d\n", *p)

	// set i through the pointer
	*p = 21

	fmt.Printf("The new value of i:%d\n", i)

	// point to j
	p = &j

	// divide j through the pointer
	*p = *p / 37

	// see the new value of j
	fmt.Printf("The new value of i:%d\n", j)
}

func TestPassByValue(t *testing.T) {
	x := 10
	fmt.Println(x)

	//calling a function with an argument like this, will pass the value, not the reference
	setToZero(x)
	fmt.Println(x)

	//calling a function with an argument like this, will pass the reference, so the value of the pointer can be
	//modified
	setToZeroThroughPointer(&x)
	fmt.Println(x)

	array := [3]int{1, 2, 3}

	for i := range array {
		array[i] = 0
	}

	fmt.Println(array)
}

//setToZero will not modify the original x, only the value received
func setToZero(x int) {
	x = 0
}

//setToZeroThroughPointer will modify the value of the pointer received as parameter
func setToZeroThroughPointer(x *int) {
	*x = 0
}

//func type
type HelloFunc func(string)

func SayHello(to string) {
	fmt.Printf("Hello, %s!\n", to)
}

func SayGreetings(to string) {
	fmt.Printf("Greetings, %s!\n", to)
}

func SayHi(helloFunc HelloFunc, to string) {
	helloFunc(to)
}

func TestPassFuncAsParam(t *testing.T) {
	SayHi(SayHello, "there")
	SayHi(SayGreetings, "there")
}
