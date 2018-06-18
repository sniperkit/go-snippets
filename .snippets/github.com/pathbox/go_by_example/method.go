package main

import (
	"fmt"
)

type rect struct {
	width, height int
}

func (r *rect) area() int { //地址拷贝
	return r.width * r.height
}

func (r rect) perim() int { // 值拷贝
	return 2*r.width + 2*r.height
}

func (r *rect) padd(a int) int { //方法不能有多个返回值, 函数可以有多个返回值
	fmt.Println("width: ", r.width)
	r.width += a
	return r.width
}

func (r rect) iadd(a int) int {
	fmt.Println("width: ", r.width)
	r.width += a
	return r.width // r.width 只是在这个范围中 被修改了, 实际struct中的 width是没有被修改的
}

func main() {
	r := rect{width: 10, height: 10}

	fmt.Println("area: ", r.area())
	fmt.Println("perim:", r.perim())
	rp := &r
	fmt.Println("area: ", rp.area())
	fmt.Println("perim:", rp.perim())

	a := 10
	fmt.Println("iadd: ", r.iadd(a))
	fmt.Println("padd: ", r.padd(a))
	fmt.Println("padd: ", r.padd(a))

}
