package main

import "fmt"

var cnt = 0

func main() {
	fmt.Println(cnt)
	cnt++
	fmt.Println(cnt)
	info()
	fmt.Println(cnt)
}

func info() {
	fmt.Println(cnt)
	cnt += 10
	fmt.Println(cnt)
}
