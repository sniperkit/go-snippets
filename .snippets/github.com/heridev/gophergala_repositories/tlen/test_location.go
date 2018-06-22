package main

type TestLocation struct {
	BaseLocation
}

func NewTestLocation(title, desc string) *TestLocation {
	return &TestLocation{BaseLocation{ButtonTitle: title, Description: desc}}
}
