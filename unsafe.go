package main

import (
	"fmt"
	"unsafe"
)

const (
	FOOSIZE = 16
)

type foo struct {
	id    uint64
	index uint32
	align uint32
}

func main() {
	fmt.Println("unsafe testing")

	var fooBuf [FOOSIZE]byte
	fo := (*foo)(unsafe.Pointer(&fooBuf[0]))
	fo.id = 123
	fo.index = 234
	fo.align = 456

	fmt.Printf("fooBuf = %+v\n", fooBuf)
	fmt.Printf("foo struct = %+v\n", fo)
}
