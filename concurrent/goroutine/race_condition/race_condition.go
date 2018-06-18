// 竞争状态
package main

import (
	"sync"
	"fmt"
	"runtime"
)

var (
	counter int
	wait    sync.WaitGroup
)

func main() {
	// 等待2个 goroutine
	wait.Add(2)

	// 创建2个 goroutine
	go incCounter(1)
	go incCounter(2)

	wait.Wait()
	fmt.Println("Final Counter:", counter)
}

// incCounter
func incCounter(id int) {
	// 在函数退出时调用 Done 来通知 main 函数工作已经完成
	defer wait.Done()

	for count := 0; count < 2; count++ {
		// 捕获 counter 的值
		value := counter

		// 当前 goroutine 从线程退出，并放回到队列「强制调度器切换2个 goroutine，以便让竞争状态的效果明显」
		runtime.Gosched()

		// 增加本地 value 变量的值
		value++
		// 将该值保存会 counter
		counter = value
	}
}
