//https://segmentfault.com/q/1010000000198984

package main

import "fmt"

type DeptModeFull interface {
	Name() string
	SetName(name string)
	Relocate(building string, floor uint8)
}

type DeptModeA interface {
	Name() string
	SetName(name string)
}

type DeptModeB interface {
	Relocate(building string, floor uint8)
}

type Dept struct {
	name     string
	building string
	floor    uint8
	Key      string
}

func (self Dept) Name() string {
	return self.name
}

func (self Dept) SetName(name string) {
	self.name = name
}

func (self *Dept) Relocate(building string, floor uint8) {
	self.building = building
	self.floor = floor
}

func main() {

	dept1 :=
		Dept{
			name:     "MySohu",
			building: "Internet",
			floor:    7}
	switch v := interface{}(dept1).(type) {
	case DeptModeFull:
		fmt.Printf("The dept1 is a DeptModeFull.\n")
	case DeptModeB:
		fmt.Printf("The dept1 is a DeptModeB.\n")
	case DeptModeA:
		fmt.Printf("The dept1 is a DeptModeA.\n")
	default:
		fmt.Printf("The type of dept1 is %v\n", v)
	}
	deptPtr1 := &dept1
	if _, ok := interface{}(deptPtr1).(DeptModeFull); ok {
		fmt.Printf("The deptPtr1 is a DeptModeFull.\n")
	}
	if _, ok := interface{}(deptPtr1).(DeptModeA); ok {
		fmt.Printf("The deptPtr1 is a DeptModeA.\n")
	}
	if _, ok := interface{}(deptPtr1).(DeptModeB); ok {
		fmt.Printf("The deptPtr1 is a DeptModeB.\n")
	}
}
