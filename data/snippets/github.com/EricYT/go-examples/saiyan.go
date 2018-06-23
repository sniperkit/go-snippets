package main

import (
  "fmt"
)

type Person struct {
  Name string
}

func (person *Person) Introduce() {
  fmt.Printf("Hi, I'm %s\n", person.Name)
}

type Saiyan struct {
  *Person
  Power int
  Name string
}

func main() {
  goku := &Saiyan{
    Person : &Person{Name:"yutao"},
    Power  : 1,
    Name : "eric",
  }

  goku.Introduce()
  goku.Person.Introduce()
  fmt.Println(goku.Power)
  fmt.Println(goku.Name)
  fmt.Println(goku.Person.Name)
}
