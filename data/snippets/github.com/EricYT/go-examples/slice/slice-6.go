package main

import "fmt"

func main() {
	s := make([]int64, 2)

  s[0] = int64(1)
  s[1] = int64(2)

	fmt.Printf("address of slice %p len %d cap %d\n", &s, len(s), cap(s))

	for i := 0; i < 10000; i++ {
    s = append(s, int64(i))
	}

	fmt.Printf("address of slice %p len %d cap %d\n", &s, len(s), cap(s))
	// Not enough space, twice the captial

  s1 := []int{1, 2}
	fmt.Printf("address of slice %p len %d cap %d\n", &s1, len(s1), cap(s1))
  s1 = append(s1, 3)
  s1 = append(s1, 3)
  s1 = append(s1, 3)
	fmt.Printf("address of slice %p len %d cap %d\n", &s1, len(s1), cap(s1))

}
