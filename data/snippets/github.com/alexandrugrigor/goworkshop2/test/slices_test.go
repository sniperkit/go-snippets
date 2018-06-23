package test

import (
	"fmt"
	"testing"
)

type SliceType []string

func TestSlices(t *testing.T) {

	//1. How to declare?
	var a []int
	fmt.Printf("\n1. a = %v\n", a)

	// 2. How to initialize?
	slice := []string{"a", "b", "c", "d"}
	slice2 := make([]int, 3)
	fmt.Printf("\n2. slice = %v\n", slice)
	fmt.Printf("2. slice2 = %v\n", slice2)

	// 3. How to get the size?
	fmt.Printf("\n3. slice size = %d\n", len(slice))
	fmt.Printf("3. slice2 size = %d\n", len(slice2))

	// 4. How to get an element
	fmt.Printf("\n4. slice[2] = %s\n", slice[2])

	// 5. How to add an element?
	slice = append(slice, "e")
	fmt.Printf("5. slice[4] = %s\n", slice[4])

	// 6. How to traverse?
	fmt.Printf("\n6. traversing the slice:\n")
	for index, value := range slice {
		fmt.Printf("\tslice[%d]=%s\n", index, value)
	}

	// 7. How to remove an element?
	var aSlice SliceType = []string{"a", "b", "c", "d"}
	aSlice.remove("c")
	fmt.Printf("\n7. removing element 'c' from aSlice:\n")
	fmt.Printf("aSlice = %v\n", aSlice)

	// 8. How to update an element?
	aSlice.update("b", "bb")
	fmt.Printf("\n8. updating element 'b' to 'bb' of aSlice:\n")
	fmt.Printf("aSlice = %v\n", aSlice)
}

func (p *SliceType) remove(element string) error {
	var err = fmt.Errorf("could not find element %s", element)
	var updatedSlice SliceType
	for _, value := range *p {
		if value == element {
			err = nil
		} else {
			updatedSlice = append(updatedSlice, value)
		}
	}
	if err == nil {
		*p = updatedSlice
	}
	return err
}

func (p *SliceType) update(oldElement string, newElement string) error {
	var err = fmt.Errorf("could not find element %s", oldElement)
	var updatedSlice SliceType
	for _, value := range *p {
		if value == oldElement {
			updatedSlice = append(updatedSlice, newElement)
			err = nil
		} else {
			updatedSlice = append(updatedSlice, value)
		}
	}
	if err == nil {
		*p = updatedSlice
	}
	return err
}
