package main

import (
	"runtime"
	"sync"
	"fmt"
)
func main() {

	memConsumed := func() uint64 {
		runtime.GC()
		var s runtime.MemStats
		runtime.ReadMemStats(&s)
		return s.Sys
		// 通过runtime模块GC然后获取内存状况
	}

	var c <-chan interface{} 
	var wg sync.WaitGroup
	noop := func() { wg.Done(); <-c} // 这样go协程在主线程退出之前不会退出

	const numGoroutines = 1e4 // 创建go协程数量
	wg.Add(numGoroutines)
	before := memConsumed()
	for i := numGoroutines; i>0; i--{
		go noop()
	}
	wg.Wait()
	after := memConsumed()
	fmt.Printf("%.3f kb",float64(after-before)/numGoroutines/1000)
}
