package main

import (
	"fmt"
)

const (
	FirstName = "eric"
	LastName  = "yu"
)

type Person struct {
	FirstName, LastName string
}

func (p *Person) Introduce() {
	fmt.Println(" Persion output")
}

func (p *Person) Show() {
	fmt.Println("FirstName: ", p.FirstName)
	fmt.Println("LastName: ", p.LastName)
}

func MakePersion(firstName string, lastName string) *Person {
	return &Person{
		FirstName: firstName,
		LastName:  lastName,
	}
}

func MakePersion1(firstName string, lastName string) Person {
	return Person{
		FirstName: firstName,
		LastName:  lastName,
	}
}

// Saiyan
type Saiyan struct {
	*Person
	Power int
}

func (s *Saiyan) Introduce() {
	fmt.Println(" Saiyan output")
	s.Show()
}

func MakeSaiyan(firstName string, lastName string, power int) Saiyan {
	return Saiyan{
		Person: &Person{
			FirstName: firstName,
			LastName:  lastName,
		},
		Power: power,
	}
}

// func
type Add func(a int, b int) int

func process(add Add) int {
	return add(1, 2)
}

func main() {
	per1 := MakePersion(FirstName, LastName)
	per2 := MakePersion1(FirstName, LastName)

	per1.Show()
	per2.Show()

	saiyan := MakeSaiyan(FirstName, LastName, 100)
	saiyan.Person.Introduce()
	saiyan.Introduce()

	// func
	fmt.Println(process(func(a int, b int) int { return a + b }))
}
