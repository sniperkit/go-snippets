package main

import (
	"fmt"
	"unsafe"
)

func main() {
	var i int32 = 1
	p := unsafe.Pointer(&i)
	if *(*byte)(p) == 1 {
		fmt.Println("LittleEndian")
	} else {
		fmt.Println("BigEndian")
	}
}
