package main

import (
	"log"
	"time"
	"github.com/jiangew/hancock/concurrent/worker"
	"sync"
)

var names = []string{
	"Lily",
	"JamesiWorks",
	"LaoLuo",
	"Steve",
	"Jackson",
	"Mary",
}

type namePrinter struct {
	name string
}

// Task 实现 Worker 接口
func (m *namePrinter) Task() {
	log.Println(m.name)
	time.Sleep(time.Second)
}

func main() {
	// 使用2个 goroutine 来创建工作池
	p := worker.New(2)

	var wg sync.WaitGroup
	wg.Add(100 * len(names))

	for i := 0; i < 100; i++ {
		for _, name := range names {
			// 创建 namePrinter 并提供指定的名字
			np := namePrinter{
				name: name,
			}

			go func() {
				// 将任务提交并执行；当Run返回时我们就知道任务已经处理完成
				p.Run(&np)
				wg.Done()
			}()
		}
	}

	wg.Wait()

	// 让工作池停止工作，等待所有现有的工作完成
	p.Shutdown()
}
