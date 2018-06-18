package main

import (
	"fmt"
)

type HashSet struct {
	m map[interface{}]bool
}

func (set *HashSet) Add(e interface{}) bool {
	if !set.m[e] {
		set.m[e] = true
		return true
	}
	return false
}

func (set *HashSet) Remove(e interface{}) {
	delete(set.m, e)
}

func (set *HashSet) Clear(e interface{}) {
	set.m = make(map[interface{}]bool)
}

func (set *HashSet) Len() int {
	return len(set.m)
}

func (set *HashSet) Contains(e interface{}) bool {
	return set.m[e]
}

func (set *HashSet) Same(other *HashSet) bool {
	if other == nil {
		return false
	}
	if set.Len() != other.Len() {
		return false
	}
	for key := range set.m {
		if !other.Contains(key) {
			return false
		}
	}
	return true
}

func (set *HashSet) Elements() []interface{} {
	initialLen := len(set.m)
	snapshot := make([]interface{}, initialLen)
	actualLen := 0
	for key := range set.m {
		if actualLen < initialLen {
			snapshot[actualLen] = key
		} else {
			snapshot = append(snapshot, key)
		}
		actualLen++
	}
	if actualLen < initialLen {
		snapshot = snapshot[:actualLen]
	}
	return snapshot
}

type Set interface {
	Add(e interface{}) bool
	Remove(e interface{}) bool
	Clear()
	Contains(e interface{}) bool
	Len() int
	Same(other Set) bool
}

func main() {
	set := &HashSet{make(map[interface{}]bool)}
	tmp := set.Add("a")
	fmt.Println(tmp)
	fmt.Println(set.m)
}
