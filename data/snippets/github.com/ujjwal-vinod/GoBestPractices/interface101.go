// 1) Functions with a pointer argument must take a pointer.
// 2) While methods with pointer receivers take either a value 
// or a pointer as the receiver when they are called. That is, 
// as a convenience, Go interprets the statement with ampersand 
// wherever required.

package main

import (
	"fmt"
)

type MyName struct {
	firstname string
	lastname  string
}

func NewMyName(f string, l string) *MyName {
	temp := new(MyName)
	temp.firstname = f
	temp.lastname = l
	return temp
}

type Setter interface {
	Setfirstname(string)
	Setlastname(string)
	Getfirstname() string
	Getlastname() string
}

func (m *MyName) Setfirstname(fn string) {
	m.firstname = fn
}

func (m *MyName) Setlastname(ln string) {
	m.lastname = ln
}

func (m *MyName) Getfirstname() string {
	return m.firstname
}

func (m *MyName) Getlastname() string {
	return m.lastname
}

func main() {
	mn := NewMyName("Ujjwal", "Vinod")
	fmt.Println(mn.Getfirstname())
	fmt.Println(mn.Getlastname())
	mn.Setlastname("Ujjwal")
	mn.Setfirstname("Arihant")
	fmt.Println(mn.Getfirstname())
	fmt.Println(mn.Getlastname())
}
