/**
 * 并发执行的数据同步问题
 * 通过加锁 Mutex.Clock 解决并发过程中的数据同步问题
 */
package main

import (
	"fmt"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

func main() {
	// 设置真正意义上的并发
	runtime.GOMAXPROCS(runtime.NumCPU())

	wg := sync.WaitGroup{}
	wg.Add(5)

	// 生成随机种子
	rand.Seed(time.Now().Unix())

	// 并发5个goroutine卖票
	for i := 0; i < 5; i++ {
		go sellTickets(&wg, i)
	}

	wg.Wait()

	// 退出时打印还有多少余票
	fmt.Println(totalTickets, "done")
}

var totalTickets int32 = 10
var mutex = &sync.Mutex{}

func sellTickets(wg *sync.WaitGroup, i int) {
	for totalTickets > 0 {
		mutex.Lock()
		if totalTickets > 0 {
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			totalTickets--
			fmt.Println("id:", i, " tickets:", totalTickets)
		}
		mutex.Unlock()
	}
	wg.Done()
}
