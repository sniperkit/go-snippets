package main

import (
	"fmt"
	"net/http"
	"runtime"
	_ "time"
)

func main() {
	header := make(http.Header)
	header.Set("X-Foo", "bar")
	header.Set("X-Fength", "12333")
	fmt.Printf("header is:%+v", header)
	runtime.LockOSThread()

	type tmp interface{}
	t := tmp{}
	switch t.(type) {
	case string:
		fmt.Println("type is a string")
	default:
		fmt.Println("unknow type")
	}

	return
}
