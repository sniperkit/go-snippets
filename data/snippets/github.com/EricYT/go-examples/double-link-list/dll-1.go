package main

import (
	"fmt"
)

type DList struct {
	id int

	prev *DList
	next *DList
}

func NewDList(i int) *DList {
	dl := &DList{id: i}
	dl.initialize()
	return dl
}

func (dl DList) Id() int {
	fmt.Println("node id:", dl.id)
	return dl.id
}
