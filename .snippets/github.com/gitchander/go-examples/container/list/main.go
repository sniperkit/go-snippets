package main

import (
	"container/list"
	"fmt"
)

func main() {

	var x list.List

	for i := 0; i < 10; i++ {
		x.PushBack(i)
	}

	for e := x.Front(); e != nil; e = e.Next() {
		fmt.Printf("%d ", e.Value.(int))
	}
	fmt.Println()

	for e := x.Back(); e != nil; e = e.Prev() {
		fmt.Printf("%d ", e.Value.(int))
	}
	fmt.Println()
}
