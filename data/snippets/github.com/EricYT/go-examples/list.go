package main

import (
	"container/list"
	"fmt"
)

func main() {
	var l list.List
	l.PushBack(1)
	l.PushBack(2)
	l.PushBack(3)
	l.PushBack(4)

	for value := l.Front(); value != nil; value = value.Next() {
		fmt.Println(value.Value.(int))
	}
}
