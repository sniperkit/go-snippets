package main

import (
	"fmt"
	"sync/atomic"
)

func main() {
	var m atomic.Value
	m.Store("hello,world")
	fmt.Println("atomic value is ", m.Load().(string))
}
