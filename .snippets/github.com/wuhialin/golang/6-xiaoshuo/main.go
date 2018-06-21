package main

import (
	_ "./common"
	_ "./crawler"
	//_ "./data"
	"runtime"
	"time"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.Gosched()
	for range time.Tick(time.Second) {
	}
}
