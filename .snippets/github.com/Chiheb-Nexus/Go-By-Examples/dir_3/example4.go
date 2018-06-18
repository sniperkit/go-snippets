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

func returnMaxLength(tab []string) int { // Private function
	max_length := 0
	for _, value := range tab {
		if len(value) > max_length {
			max_length = len(value)
		}
	}
	return max_length
}

func GroupByLength(tab []string) [][]string { // Public function
	max_length := returnMaxLength(tab)
	tab_new := make([][]string, max_length)
	for _, name := range tab {
		tab_new[len(name)-1] = append(tab_new[len(name)-1], name)
	}
	return tab_new
}

func main() {
	fmt.Println(GroupByLength(names))
}
