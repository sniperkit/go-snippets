package main

import (
	"fmt"
)

type List struct {
	index int

	prev *List
	next *List
}

func (l List) Display() {
	fmt.Println("List index:", l.index)
}

func NewList(i int) *List {
	l := &List{index: i}
	l.prev = l
	l.next = l
	return l
}

func ListInsert(list, prev, next *List) {
	next.prev = list
	list.next = next
	list.prev = prev
	prev.next = list
}

func ListSplice(prev, next *List) {
	prev.next = next
	next.prev = prev
}

func ListRemove(head *List) {
	ListSplice(head.prev, head.next)
	head.next = nil
	head.prev = nil
}

func ListPrepend(head, list *List) {
	ListInsert(head, list, list.next)
}

func ListAppend(head, list *List) {
	head.next.prev = list
	list.next = head.next
	head.next = list
	list.prev = head
}

func ListPush(head, list *List) *List {
	head.prev.next = list
	list.prev = head.prev
	head.prev = list
	list.next = head
	return list
}

func Foreach(head *List) {
	head.Display()
	for i := head.next; i != head; i = i.next {
		i.Display()
	}
}

func main() {
	l1 := NewList(1)
	fmt.Printf("l1 = %+v\n", l1)
	l2 := NewList(2)
	fmt.Printf("l2 = %+v\n", l2)
	l3 := NewList(3)
	fmt.Printf("l3 = %+v\n", l3)

	//ListPrepend(l2, l1)
	//ListAppend(l1, l2)
	l1 = ListPush(l1, l2)
	Foreach(l1)
	fmt.Printf("l1 + l2 = %+v\n", l1)

	//ListPrepend(l3, l1)
	//ListAppend(l1, l3)
	l1 = ListPush(l1, l3)
	Foreach(l1)
	fmt.Printf("l1 + l3 = %+v\n", l1)
}
