package main

import (
	"fmt"
	"time"
)

func main() {
	t1 := task1()

	ch1 := make(chan T1)
	go func() {
		time.Sleep(time.Second * 2)

		// once task2 is done, inform ch1
		ch1 <- task2(t1)
	}()

	ch2 := make(chan T2)
	go func() {
		time.Sleep(time.Second * 3)

		// once task3 is done, inform ch2
		ch2 <- task3(t1)
	}()

	var t1Received T1
	var t2Received T2

	const loopIntervalSeconds = time.Second
	for {
		select {
		// send `ch1` to `t1Received`
		case t1Received = <-ch1:
			fmt.Println(t1Received, "channel received")

		// send `ch2` to `t2Received`
		case t2Received = <-ch2:
			fmt.Println(t2Received, "channel received")

		default:
			fmt.Println("Waiting for channel..")
		}
		if (t1Received != T1{} && t2Received != T2{}) {
			break
		}
		time.Sleep(loopIntervalSeconds)
	}

	// task4 will be executed as soon as
	// `t1Received` & `t2Received` arrive
	task4(t1Received, t2Received)
}

type T1 struct {
	name string
}

type T2 struct {
	name string
}

func task1() T1 {
	fmt.Println("Task 1 done")
	return T1{"Task1"}
}

func task2(t1 T1) T1 {
	fmt.Println("Task 2 done")
	return t1
}

func task3(t1 T1) T2 {
	fmt.Println("Task 3 done")
	return T2{"Task2"}
}

func task4(t1 T1, t2 T2) {
	fmt.Println("Task 4 done")
}
