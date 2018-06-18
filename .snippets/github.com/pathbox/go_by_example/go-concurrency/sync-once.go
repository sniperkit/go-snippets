package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var num int
	sign := make(chan bool)
	var once sync.Once
	f := func(ii int) func() {
		return func() {
			num = (num + ii*2)
			sign <- true
		}
	}
	for i := 0; i < 3; i++ {
		fi := f(i + 1)
		go once.Do(fi)
		for j := 0; j < 3; j++ {
			select {
			case <-sign:
				fmt.Println("Receive a signal")
			case <-time.After(100 * time.Millisecond): // 超时设置
				fmt.Println("Timeout!")
			}
		}
		fmt.Printf("Num: %d\n", num)
	}
}

// Receive a signal
// Timeout!
// Timeout!
// Num: 2
// Timeout!
// Timeout!
// Timeout!
// Num: 2
// Timeout!
// Timeout!
// Timeout!
// Num: 2
// fi just run one time

// sync.Once类型的典型应用场景就是执行仅需要执行一次的任务。例如：数据块连接池初始化任务。又例如：一些需要持续运行的实时监测任务等等
