package main

import (
	"fmt"
	"time"
)

func main() {
	relax(3)
}

func relax(sec int) {
	fmt.Println("i take a nap")
	time.Sleep(time.Duration(sec) * time.Second)
	fmt.Println("was nice to nap for ", sec, " seconds")
}
