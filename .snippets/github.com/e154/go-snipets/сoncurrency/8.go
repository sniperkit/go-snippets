package main

import (
	"sync"
	"time"
	"syscall"
	"os"
	"encoding/json"
)

var (
	Stdout = os.NewFile(uintptr(syscall.Stdout), "/dev/stdout")
)

func NewModel() (model *Model) {
	model = &Model{
		count: make(map[int]int),
		push: make(chan int),
		get: make(chan map[int]int),
	}

	go model.loop()

	return
}

type Model struct {
	sync.RWMutex
	count map[int]int
	get   chan map[int]int
	push  chan int
}

func (m *Model) loop() {

	for {
		select {
		case num := <- m.push:
			m.count[num] = num
		case m.get <- m.count:
		}

		time.Sleep(20)
	}
}

func (m *Model) Push(num int) {
	m.push <- num
}

func (m *Model) GetCounter() map[int]int {
	return <- m.get
}

func push(model *Model) {
	for i:=0;i<20;i++ {
		model.Push(i)
	}
}

func get(model *Model)  {
	v := model.GetCounter()
	j, _ := json.Marshal(v)
	Stdout.Write(j)
	//log.Println(v)
}

func main() {
	model := NewModel()
	for i:=0;i<2;i++ {
		go push(model)
	}

	for i:=0;i<10;i++ {
		get(model)
	}

	//for {
	//	log.Println(model.GetCounter())
	//
	//}
}