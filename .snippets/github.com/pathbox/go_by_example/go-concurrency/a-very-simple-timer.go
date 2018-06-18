package main

import (
	"fmt"
	"time"
)

func main() {
	// go func(){
	// 	var timer *time.Timer
	// 	for {
	// 		// 省略若干语句
	// 		case <-func() <-chan time.Timer{
	// 			if timer == nil {
	// 				timer = time.NewTimer(time.Millisecond)
	// 			}else{
	// 				timer.Reset(time.Millisecond)
	// 			}
	// 			return time.C
	// 		}():
	// 		fmt.Println("Timeout.")
	// 		ok = false
	// 		break
	// 	}
	// }()

	var t *time.Timer
	f := func() {
		fmt.Printf("Expiration time: %v\n", time.Now())
		fmt.Printf("C's len: %d\n", len(t.C))
	}
	t = time.AfterFunc(1*time.Second, f) // 异步的方式在到期事件来临的那一刻执行定义的函数。不是阻塞的
	time.Sleep(2 * time.Second)
}
