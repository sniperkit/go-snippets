package main

import (
	"fmt"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}
	n := max(nums)
	fmt.Println(n)
	a := avg(nums)
	fmt.Println(a)
	ary := rev_array(nums)
	fmt.Println(ary)

}

func max(nums []int) int {
	m := nums[0]

	for _, num := range nums {
		if num > m {
			m = num
		}
	}
	return m
}

func avg(nums []int) int {
	n := len(nums)
	m := 0
	for _, num := range nums {
		m = m + num
	}
	avg := m / n
	return avg
}

func rev_array(nums []int) []int {
	l := len(nums)
	for i, _ := range nums {
		t := nums[i]
		nums[i] = nums[l-1-i]
		nums[l-1-i] = t
	}
	return nums
}
