package main

import (
	"fmt"
	"sort"
)

// 12. Каким образом отсортировать слайс целых чисел (var arr[]int)?

// Ответы:
// sort.Ints()

func main() {

	i := []int{2,4,1,5}

	sort.Ints(i)

	fmt.Println(i)
}
