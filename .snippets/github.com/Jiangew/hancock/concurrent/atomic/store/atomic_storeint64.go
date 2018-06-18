// 使用 atomic 原子操作保证对数值类型的安全访问
package main

import (
	"sync"
	"time"
	"fmt"
	"sync/atomic"
)

var (
	// 通知正在执行的 goroutine 停止工作的标志
	shutdown int64
	wait     sync.WaitGroup
)

func main() {
	wait.Add(2)

	go doWork("A")
	go doWork("B")

	// 给定 goroutine 执行时间
	time.Sleep(1 * time.Second)

	fmt.Println("Shutdown Now")

	// 该停止工作了，安全的设置 shutdown 标志
	atomic.StoreInt64(&shutdown, 1)

	wait.Wait()
}

// doWork 检测 shutdown 标志来决定是否提前终止
func doWork(name string) {
	defer wait.Done()

	for {
		fmt.Printf("Done %s Work\n", name)
		time.Sleep(250 * time.Millisecond)

		// 要停止工作了吗？
		if atomic.LoadInt64(&shutdown) == 1 {
			fmt.Printf("Shutting %s Down\n", name)
			break
		}
	}
}
