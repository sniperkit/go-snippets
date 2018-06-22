package main

import "fmt"

func main() {
	// var a []int
	// var b = []string{"a", "b"}
	// c := []string{"a", "b"}

	stringSlice1 := []string{"1", "2", "3", "4"}
	fmt.Println(stringSlice1[:2])
	fmt.Println(stringSlice1[2:])

	stringSlice2 := make([]string, 5, 10)
	copy(stringSlice2, stringSlice1)
	fmt.Println(stringSlice2[0])

	stringSlice2 = append(stringSlice2, "5", "6")
	fmt.Println(stringSlice2)
	fmt.Println(stringSlice2[5])
	fmt.Println(stringSlice2[6])
	stringSlice2 = append(stringSlice2, []string{"7", "8"}...)
	fmt.Println(stringSlice2)
	fmt.Println(stringSlice2[7])
	fmt.Println(stringSlice2[8])
}
