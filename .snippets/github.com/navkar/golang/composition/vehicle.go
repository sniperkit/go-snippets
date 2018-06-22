package main

import (
  "fmt"
)

type Vehicle struct {
  Wheels int64
  IsElectric bool
  EngineType string
}

type Car struct {
  *Vehicle
  Price float64
}

func (v *Vehicle) CreateVehicle() {
  v.Wheels = 4
  v.IsElectric = false
  v.EngineType = "V4"
}

// Go supports composition, which is the act of including one structure into another.
// In some languages, this is called a trait or a mixin.
// Languages that don’t have an explicit composition mechanism can always do it the long way.
func main() {
  v := new (Vehicle)
  v.CreateVehicle()

  c := new (Car)
  c.Price = 239.49
  c.Vehicle = v

  // Important : use this c.Vehicle.EngineType OR c.EngineType
  fmt.Println("The vehicle has " + c.Vehicle.EngineType + " engine")
  fmt.Println("The price of the car is " + fmt.Sprintf("%.2f", c.Price))

}
