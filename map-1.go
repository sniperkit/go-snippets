package main

import "fmt"
import _ "reflect"

type Node struct {
	Next  *Node
	Value interface{}
}

func main() {
	var first *Node

	/* This is a double linked list */
	var a Node
	var b Node

	a = Node{
		Next:  &b,
		Value: 1,
	}
	b = Node{
		Next:  &a,
		Value: 2,
	}

	first = &b

	visited := make(map[*Node]bool)
	for n := first; n != nil; n = n.Next {
		// fmt.Println("TypeOf ", reflect.TypeOf(visited[n]))
		if visited[n] {
			fmt.Println("cycle deletced")
			break
		}
		visited[n] = true
		fmt.Println("Value is ", n.Value)
	}
}
