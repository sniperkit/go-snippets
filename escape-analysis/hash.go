package main

import (
	"fmt"
	"hash/fnv"
)

func hashIt(in string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(in))
	out := h.Sum64()
	return out
}

func main() {
	s := "hello"
	fmt.Printf("The FNV6a hash of '%v' is '%v'\n", s, hashIt(s))
}
