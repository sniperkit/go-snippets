package main

import (
	"fmt"
	"time"
)

func main() {
	go func() {
		var e int
		ok := true
		for {
			select {
			case e, ok = <-ch11:
			// 省略若干代码
			case ok = <-func() chan bool {
				timeout := make(chan bool, 1)
				go func() {
					time.Sleep(time.Millisecond)
					timeout <- false
				}()
				return timeout
			}():
				fmt.Println("Timeout")
				break
			}
			if !ok {
				break
			}
		}
	}()
}

// 在运行时系统开始执行select语句的时候，会先对它所有的case中的元素表达式和通道表达式进行求值。
// 这样才使得在运行时系统选择要执行case之前先制造出一个可用的超时触发器成为了可能。更具体地讲，在
// 这些case被选择之前，第二个case后的接收语句会由下面这行代码代替： ok = <- timeout

// 运行时系统在选择select语句的case的时候，只要case有多个，它就肯定不会为某一个case而等待。只有当
// 所有的case后的发送语句或接收语句都无法被立即执行的时候，它才会阻塞住当前的GOroutine。当然，前提是
// 没有default case。在等待期间，只要发现有某一个case后的语句可以立即被执行，那么运行时系统就会立即执行这个case
// 在本例中，当无法立即从ch11通道中接收元素值的时候，运行时系统会随机判断是否可以立即接收timeout通道中的元素值。
// 因此，一旦第一个case中的接收操作无法在1毫秒之内完成，我们给定超时子流程就会被执行
