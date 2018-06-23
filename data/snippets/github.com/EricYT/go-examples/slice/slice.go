package main

import "fmt"

func main() {
	s1 := []int{1, 2, 3, 4}
	s2 := s1[:]
	s1[0] = 9
	fmt.Println(s1)
	fmt.Println(s2)
	s2 = append(s2, 5, 6, 7, 8)
	fmt.Println(s2)
	// s1 and s2 point to the same array

	s3 := make([]int, 3)
	copy(s3, s2)
	fmt.Println(s3)
	s3[0] = 1
	fmt.Println(s3)
	fmt.Println(s2)
	// s3 and s2 have diffirence place to store their elements
}
