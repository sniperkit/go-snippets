/**
 * 无缓冲通道：使用通道，在 goroutine 之间进行数据同步
 * 4个 goroutine 间的接力比赛
 */
package main

import (
	"sync"
	"fmt"
	"time"
)

var wg sync.WaitGroup

func main() {
	// 创建无缓冲通道
	baton := make(chan int)

	// 为最后一位跑步者将计数加1
	wg.Add(1)

	// 第一位跑步者持有接力棒
	go Runner(baton)

	// 开始比赛
	baton <- 1

	wg.Wait()
}

// Runner 模拟接力赛中的一位跑步者
func Runner(baton chan int) {
	var newRunner int

	// 等待接力棒
	runner := <-baton

	// 开始绕着跑道跑步
	fmt.Printf("Runner %d Running With Baton\n", runner)

	// 创建下一位跑步者
	if runner != 4 {
		newRunner = runner + 1
		fmt.Printf("Runner %d To The Line\n", newRunner)
		go Runner(baton)
	}

	// 围绕跑道跑
	time.Sleep(100 * time.Millisecond)

	// 比赛结束了吗？
	if runner == 4 {
		fmt.Printf("Runner %d Finished, Race Over\n", runner)
		wg.Done()
		return
	}

	fmt.Printf("Runner %d Exchange With Runner %d\n", runner, newRunner)
	baton <- newRunner
}
