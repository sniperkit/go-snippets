package channel

/*
 *
 * go happens-before常用的三原则是：
 *
 * 1. 对于不带缓冲区的channel，对其写happens-before对其读
 * 2. 对于带缓冲区的channel,对其读happens-before对其写
 * 3. 对于不带缓冲的channel的接收操作 happens-before 相应channel的发送操作完成
 */

import "testing"

func TestHappensBefore_1(t *testing.T) {
	var c = make(chan int, 10)
	var a string
	f := func() {
		a = "hello, world"
		c <- 0
	}

	go f()
	<-c
	print(a)
}

func TestHappensBefore_2(t *testing.T) {
	var c = make(chan int)
	var a string
	f := func() {
		a = "hello, world"
		<-c
	}

	go f()
	c <- 0
}

func TestHappensBefore_3(t *testing.T) {
	var c = make(chan int, 1)
	var a string
	f := func() {
		a = "hello, world"
		<-c
	}

	go f()
	c <- 0
	print(a)
}
