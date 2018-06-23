package main

type FamilyBunkLocation struct {
	BaseLocation
	Header      string
	Description string
}

func NewFamilyBunkLocation(header, desc string) *FamilyBunkLocation {
	return &FamilyBunkLocation{Header: header, Description: desc}
}

func (location FamilyBunkLocation) GetHeader() string {
	return "header of family bunk"
}

func (location FamilyBunkLocation) GetDescription() string {
	return "desc of fam bun"
}
