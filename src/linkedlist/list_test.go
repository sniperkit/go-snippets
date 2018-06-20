package linkedlist

import (
	"testing"
	"fmt"
)

type test struct {
	v1 int
	v2 string
}

var testList = []interface{}{0, "123", -20, "test", test{v1: 2, v2: "abc"}, -21, 50.20}

func TestList_Insert(t *testing.T) {
	list := NewList()
	for i, e := range testList {
		list.Insert(i,e)
	}
	vLen, lLen := len(testList), list.Length()
	if vLen != lLen {
		t.Error(
			"For", list,
			"expected", vLen,
			"got", lLen,
		)
	}
	for _, e := range testList {
		if !list.Contains(e) {
			t.Error(
				"For", list,
				"expected", true,
				"got", false,
			)
		}
	}
	fmt.Println(list)
}

func TestList_Delete(t *testing.T) {
	list := NewList()
	for i, e := range testList {
		list.Insert(i,e)
	}
	for _, e := range testList {
		list.Delete(e)
	}
	listLen := list.Length()
	if listLen != 0 {
		t.Error(
			"For", list,
			"expected", 0,
			"got", listLen,
		)
	}
	fmt.Println(list)
}
