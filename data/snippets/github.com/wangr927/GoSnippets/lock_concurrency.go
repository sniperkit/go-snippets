package main

import (
	"sync"
	"time"
	"fmt"
)

// SafeCounter 并发安全的计数器
type SafeCounter struct {
	v     map[string]int
	mux   sync.Mutex
}

func (c *SafeCounter) Inc(key string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.v[key]++
}

func (c *SafeCounter) Value(key string) int {
	c.mux.Lock()
	defer c.mux.Unlock()
	return c.v[key]
}

func safeCounterSample() {
  c := SafeCounter{v:make(map[string]int)}
	for i := 0; i < 100; i ++ {
		go c.Inc("demokey")
		time.Sleep(100 * time.Millisecond)
		fmt.Println(c.Value("demokey"))
	}
	fmt.Println(c.Value("demokey"))
}
