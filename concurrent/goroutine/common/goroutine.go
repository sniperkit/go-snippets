// 如何创建 goroutine，以及调度器的行为
package main

import (
	"runtime"
	"sync"
	"fmt"
)

// wg 用来等待程序完成；计数器加2，表示要等待2个goroutine
var wg sync.WaitGroup

func main() {
	// 分配2个逻辑处理器给调度器使用
	runtime.GOMAXPROCS(2)
	// 给每个核分配一个逻辑处理器
	//runtime.GOMAXPROCS(runtime.NumCPU())

	// wg 用来等待程序完成；计数器加2，表示要等待2个goroutine
	//var wg sync.WaitGroup
	wg.Add(4)

	fmt.Println("Start Goroutines")

	// 声明一个匿名函数，并创建一个 goroutine
	go func() {
		// 在函数退出时调用 Done 来通知 main 函数工作已经完成
		defer wg.Done()

		for count := 0; count < 3; count++ {
			for char := 'a'; char < 'a'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	// 声明一个匿名函数，并创建一个 goroutine
	go func() {
		// 在函数退出时调用 Done 来通知 main 函数工作已经完成
		defer wg.Done()

		for count := 0; count < 3; count++ {
			for char := 'A'; char < 'A'+26; char++ {
				fmt.Printf("%c ", char)
			}
		}
	}()

	// 创建2个goroutine
	go printPrime("A")
	go printPrime("B")

	fmt.Println("Waiting To Finish")
	wg.Wait()

	fmt.Println("\nTerminating Program")
}

// printPrime 显示5000以内的素数
func printPrime(prefix string) {
	defer wg.Done()

next:
	for outer := 2; outer < 5000; outer++ {
		for inner := 2; inner < outer; inner++ {
			if outer%inner == 0 {
				continue next
			}
		}
		fmt.Printf("%s:%d\n", prefix, outer)
	}
	fmt.Println("Completed", prefix)
}
