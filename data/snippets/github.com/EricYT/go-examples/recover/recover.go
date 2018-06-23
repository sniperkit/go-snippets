package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	fmt.Println("Test normal recover and panic")

	var timeAfter = time.After

	t := timeAfter(10 * time.Second)
	fmt.Println("time is:", t)
	<-t

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		wg.Done()
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("recover receive error:", r)
			}
		}()

		fmt.Println("into goruntine")
		//		panic("crash")

	}()

	wg.Wait()
	time.Sleep(time.Second * 2)
	fmt.Println("main process end")
	return
}
