package main

import (
	"fmt"
	"math"
)

type geometry interface { // 任何类型 包括struct 都可以调用interface中的方法。可以看出interface是类型的父类
	area() float64 // rect 和 circle的方法名相同,但是根据类型不同调用不同代码
	perim() float64
}

type rect struct {
	width, height float64
}

type circle struct {
	radius float64
}

func (r *rect) area() float64 {
	return r.width * r.height
}

func (r rect) perim() float64 {
	return 2*r.width + 2*r.height
}

func (c circle) area() float64 {
	return math.Pi * c.radius * c.radius
}
func (c circle) perim() float64 {
	return 2 * math.Pi * c.radius
}

func measure(g geometry) {
	fmt.Println("here: ", g)
	fmt.Println(g.area()) // rect 和 circle的方法名相同,但是根据类型不同调用不同代码
	fmt.Println(g.perim())
}

func main() {
	r := rect{width: 10, height: 10}
	c := circle{radius: 5}

	measure(&r)
	measure(c)
}

// 接口： 接口中集成很多方法,或嵌套接口。在使用的时候,根据传的类型的不同,调用其该类型定义的方法
// go中没有继承， 所以用接口可以模拟实现类的继承。 可以看成interface就是父类, 子类调用其中的方法时,就是调用该对应子类的方法
