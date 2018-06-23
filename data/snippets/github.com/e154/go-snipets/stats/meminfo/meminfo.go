package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"sync"
)

func bigBytes() *[]byte {
	s := make([]byte, 100000000)
	return &s
}

func main() {
	var wg sync.WaitGroup

	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)
	log.Println("Alloc",mem.Alloc / 1024)
	log.Println("TotalAlloc",mem.TotalAlloc / 1024)
	log.Println("Heap",mem.HeapAlloc / 1024)
	log.Println("Sys",mem.HeapSys / 1024)

	//for i := 0; i < 10; i++ {
	//	s := bigBytes()
	//	if s == nil {
	//		log.Println("oh noes")
	//	}
	//}
	//
	//runtime.ReadMemStats(&mem)
	//log.Println("Alloc", mem.Alloc / 1024)
	//log.Println("TotalAlloc", mem.TotalAlloc / 1024)
	//log.Println("Heap", mem.HeapAlloc / 1024)
	//log.Println("Sys", mem.HeapSys / 1024)

	wg.Add(1)
	wg.Wait()

}