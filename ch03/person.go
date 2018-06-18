// Person struct with methods of pointer receiver

package main

import (
	"fmt"
	"time"
)

// Person struct
type Person struct {
	FirstName, LastName string
	Dob                 time.Time
	Email, Location     string
}

// PrintName prints the name of the Person
func (p Person) PrintName() {
	fmt.Printf("\n%s %s\n", p.FirstName, p.LastName)
}

// PrintDetails prints the details of Person
func (p Person) PrintDetails() {
	fmt.Printf("[Date of Birth: %s, Email: %s, Location: %s ]\n", p.Dob.String(), p.Email, p.Location)
}

func main() {
	/*
		       // Declare a Person variable using var
				var p Person
				// Assign values to fields
				p.FirstName="Shiju"
				p.LastName="Varghese"
				p.Dob= time.Date(1979, time.February, 17, 0, 0, 0, 0, time.UTC)
				p.Email="shiju@email.com"
				p.Location= "Kochi"

				// Declare a Person variable and initialize values using Struct literal
				p := Person{
					"Shiju",
					"Varghese",
					time.Date(1979, time.February, 17, 0, 0, 0, 0, time.UTC),
					"shiju@email.com",
					"Kochi",
				}

	*/
	p := Person{
		FirstName: "Shiju",
		LastName:  "Varghese",
		Dob:       time.Date(1979, time.February, 17, 0, 0, 0, 0, time.UTC),
		Email:     "shiju@email.com",
		Location:  "Kochi",
	}
	p.PrintName()
	p.PrintDetails()

}
