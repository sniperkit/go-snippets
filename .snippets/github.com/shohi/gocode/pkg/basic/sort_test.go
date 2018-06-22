package basic

import (
	"fmt"
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	ints := []int{1, 2, 5, 6}
	sort.Ints(ints)
	fmt.Println(ints)
	idx := sort.Search(len(ints), func(i int) bool { return (ints[i] >= 10) })
	fmt.Println(idx)
}
