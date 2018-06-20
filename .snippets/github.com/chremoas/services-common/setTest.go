package main

import (
	"github.com/chremoas/services-common/sets"
	"fmt"
)

func main() {
	set1 := sets.NewStringSet()
	set1.Add("String1")
	set1.Add("String2")
	set1.Add("String3")

	set2 := sets.NewStringSet()
	set2.Add("String2")
	set2.Add("String3")
	set2.Add("String4")

	fmt.Println(set1.Intersection(set2))
	fmt.Println(set1.Difference(set2))
	fmt.Println(set2.Difference(set1))
}