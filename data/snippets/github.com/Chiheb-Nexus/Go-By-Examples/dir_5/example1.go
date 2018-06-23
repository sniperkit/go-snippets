package main

import (
	"fmt"
)

type Person struct {
	Name, Surname string
	Age           int
}

type Dog struct {
	Name string
	Age  int
}

// There are two reasons to use a pointer receiver. First, to avoid copying the
// value on each method call (more efficient if the value type is a large struct).
func (p *Person) SaySomeThing() {
	fmt.Printf("Heyy! I'm %s %s and my age is %d\n", p.Name, p.Surname, p.Age)
}

func (d *Dog) SaySomeThing() {
	fmt.Printf("Heyy! i'm %s the dog, my age is %d. Weird i can talk ???!\n", d.Name, d.Age)
}

func main() {
	person := &Person{Name: "Chiheb", Surname: "Nexus", Age: 28}
	dog := &Dog{Name: "Rex", Age: 3}

	person.SaySomeThing()
	dog.SaySomeThing()
}
