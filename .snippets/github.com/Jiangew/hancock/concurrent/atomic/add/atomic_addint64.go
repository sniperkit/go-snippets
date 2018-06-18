// 使用 atomic 原子操作保证对数值类型的安全访问
package main

import (
	"sync"
	"fmt"
	"sync/atomic"
	"runtime"
)

var (
	counter int64
	wg      sync.WaitGroup
)

func main() {
	wg.Add(2)

	go incAtomicCounter(1)
	go incAtomicCounter(2)

	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

func incAtomicCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 安全数值访问
		atomic.AddInt64(&counter, 1)

		// 当前 goroutine 从线程退出，并放回到队列「强制调度器切换2个 goroutine，以便让竞争状态的效果明显」
		runtime.Gosched()
	}
}
