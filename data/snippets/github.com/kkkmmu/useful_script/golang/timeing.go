package main

import (
	"log"
	"time"
)

func StartTimer(name string) func() {
	t := time.Now()
	log.Println(name, "started")
	//这里通过闭包捕获了初始状态
	return func() {
		d := time.Now().Sub(t)
		log.Println(name, "took", d)
	}
}

func FunkyFunc() {
	stop := StartTimer("FunkyFunc")
	defer stop()

	time.Sleep(1 * time.Second)
}

func main() {
	FunkyFunc()
}
