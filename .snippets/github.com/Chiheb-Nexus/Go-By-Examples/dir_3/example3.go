package main

import (
	"fmt"
)

var names = []string{
	"Katrina", "Evan", "Neil", "Adam", "Martin", "Matt",
	"Emma", "Isabella", "Emily", "Madison",
	"Ava", "Olivia", "Sophia", "Abigail",
	"Elizabeth", "Chloe", "Samantha",
	"Addison", "Natalie", "Mia", "Alexis",
}

func main() {
	/*
		Excepcted output:
		[
		[] [] [Ava Mia] [Evan Neil Adam Matt Emma] [Emily Chloe] [Martin Olivia Sophia Alexis] [Katrina Madison Abigail Addison Natalie]
		[Isabella Samantha] [Elizabeth]
		]
	*/
	max_length := 0
	for _, value := range names {
		if len(value) > max_length {
			max_length = len(value)
		}
	}

	names_new := make([][]string, max_length) // Slice of a slice of strings

	for _, name := range names {
		names_new[len(name)-1] = append(names_new[len(name)-1], name)
	}

	fmt.Println(names_new)
}
