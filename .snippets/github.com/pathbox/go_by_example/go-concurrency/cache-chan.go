package main

import (
	"fmt"
	"time"
)

// func main() {
// 	unbufChan := make(chan int)
// 	go func() {
// 		fmt.Println("Sleep a second...")
// 		time.Sleep(time.Second)
// 		num := <-unbufChan
// 		fmt.Printf("Received a integer %d \n", num)
// 	}()

// 	num := 1
// 	fmt.Printf("Send integer %d...\n", num)
// 	unbufChan <- num
// 	fmt.Println("Done.")
// }

// Send integer 1...
// Sleep a second...
// Received a integer 1
// Done.

func main() {
	bufChan := make(chan int, 1)
	go func() {
		fmt.Println("Sleep a second...")
		time.Sleep(time.Second)
		num := <-bufChan
		fmt.Printf("Received a integer %d \n", num)
	}()

	num := 1
	fmt.Printf("Send integer %d...\n", num)
	bufChan <- num
	fmt.Println("Done.")
	l := len(bufChan)
	fmt.Println(l)
}

// Send integer 1...
// Done.
// 1
// 说明buffer channel 不会阻塞(是异步的),主GOroutine执行速度比go func 的GOroutine快，没有等其执行到，就继续往下
// 执行并结束了。l 的值为1 说明在 bufChan中有一个chan还没没执行完
