package main

import (
	"fmt"
	"remote/connect"
	"runtime"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	c := connect.CreateConnect("tcp://172.16.0.79:5050")
	t1 := time.Now().UnixNano()
	for j := 0; j < runtime.NumCPU(); j++ {
		wg.Add(1)
		go func(id int) {
			str := fmt.Sprintf("%s 协程ID:%d", "yuansudong", id)
			for i := 0; i < 1000000; i++ {
				if err := c.Send(0, []byte(str)); err != nil {
					fmt.Printf("错误: %s\n", err.Error())
				}
			}
			wg.Done()
		}(j)

	}
	wg.Wait()
	t2 := time.Now().UnixNano()
	fmt.Printf("本次消耗时间%d纳秒\n", t2-t1)
}
