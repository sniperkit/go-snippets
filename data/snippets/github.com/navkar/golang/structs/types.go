package main
import (
  "fmt"
)

type Employee struct {
  FirstName string
  LastName string
  Age int
  Role string
}

// Takes the address of the employee, golang passes args by value
func Super(e *Employee) {
  e.Age += 1
  e.Role = "Passed By Reference not by value"
}

func main() {
  naveen := Employee { "naveen", "karam", 0, "ByValue"}
  // golang passes args by value
  // address is being passed by value
  Super(&naveen)
  fmt.Println(naveen)
}
