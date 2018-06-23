package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

type MyStruct struct {
	ID   string
	Data string
}

func main() {
	var bin_buf bytes.Buffer

	x := MyStruct{"1", "Hello"}
	err := binary.Write(&bin_buf, binary.BigEndian, x)
	fmt.Println("error:", err)

	var y MyStruct
	err = binary.Read(&bin_buf, binary.BigEndian, y)
	fmt.Println("error:", err)
	fmt.Printf("%+v\n", y)
}
