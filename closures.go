package main

import (
	"fmt"
)

func intSeq(a int) func() int {
	i := 0
	return func() int {
		a += 1
		i += a + 1
		return i
	}
}

func main() {
	n := intSeq(2)
	fmt.Println(n())
	fmt.Println(n()) //这一次会直接进入到闭包中执行,闭包外的 i:=0 这行代码不会被执行。闭包中的变量结果都会被保留并带入下一个闭包代码执行
	fmt.Println(n())
	fmt.Println(n())

	m := intSeq(2)
	fmt.Println(m()) //另一个新的函数开始,和前面那个闭包毫无关系
	fmt.Println(n())
}
