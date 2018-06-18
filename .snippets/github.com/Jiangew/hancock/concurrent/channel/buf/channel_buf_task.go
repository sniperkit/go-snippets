/**
 * 有缓冲通道：使用通道，在 goroutine 之间进行数据同步
 * 有缓冲的通道和固定数目的 goroutine 来处理一堆工作
 */
package main

import (
	"sync"
	"math/rand"
	"time"
	"fmt"
)

const (
	numGoroutines = 4
	taskLoad      = 10
)

var waitGroup sync.WaitGroup

func init() {
	// 初始化随机数种子
	rand.Seed(time.Now().Unix())
}

func main() {
	// 创建一个有缓冲的通道来管理工作
	tasks := make(chan string, taskLoad)

	// 启动 goroutine 来处理工作
	waitGroup.Add(numGoroutines)
	for gr := 1; gr <= numGoroutines; gr++ {
		go worker(tasks, gr)
	}

	// 增加一组要完成的工作
	for post := 1; post <= taskLoad; post++ {
		tasks <- fmt.Sprintf("Task: %d", post)
	}

	// 等待所有工作都处理完时关闭通道，以便所有的 goroutine 退出
	close(tasks)

	waitGroup.Wait()
}

// worker 作为 goroutine 启动来处理从有缓冲的通道传入的工作
func worker(tasks chan string, worker int) {
	defer waitGroup.Done()

	for {
		// 等待分配工作
		task, ok := <-tasks
		if !ok {
			fmt.Printf("Worker: %d : Shutting Down\n", worker)
			return
		}
		// 显示我们开始工作了
		fmt.Printf("Worker: %d : Started %s\n", worker, task)

		// 随机等一段时间来模拟工作
		sleep := rand.Int63n(100)
		time.Sleep(time.Duration(sleep) * time.Millisecond)

		// 显示我们完成了工作
		fmt.Printf("Worker: %d : Completed %s\n", worker, task)
	}
}
