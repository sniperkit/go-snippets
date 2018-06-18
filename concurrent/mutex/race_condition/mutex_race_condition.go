/**
 * 互斥锁：用于在代码上创建一个临时区，保证同一时间只有一个 goroutine 可以执行这个临界区代码
 * 互斥锁保证资源的同步访问
 */
package main

import (
	"sync"
	"fmt"
	"runtime"
)

var (
	counter int
	wg      sync.WaitGroup

	// 定义一段代码临界区
	mux sync.Mutex
)

func main() {
	wg.Add(2)

	go incCounter(1)
	go incCounter(2)

	wg.Wait()
	fmt.Println("Final Counter:", counter)
}

// incCounter 使用互斥锁来同步并保证安全访问
func incCounter(id int) {
	defer wg.Done()

	for count := 0; count < 2; count++ {
		// 同一时刻只允许一个 goroutine 进入临界区
		mux.Lock()
		{
			value := counter

			// 当前 goroutine 从线程退出，并放回到队列
			runtime.Gosched()

			value++
			counter = value
		}
		mux.Unlock()
	}
}
