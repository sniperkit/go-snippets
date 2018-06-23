package main

import (
	"fmt"
	"sync"
	"time"
)

type message struct {
	message string
	id      int
}

func (msg *message) ShowMessage() {
	fmt.Println("message body:", msg.message)
}

func (msg *message) ShowId() {
	fmt.Println("message id:", msg.id)
}

func SendMessage(msgq chan message, msg string, id int) {
	msgTmp := message{
		message: msg,
		id:      id,
	}
	fmt.Println("SendMessage id:", id)
	msgq <- msgTmp
}

func main() {
	//wg := new(sync.WaitGroup)
	var wg sync.WaitGroup

	msg := make(chan message, 2)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(wg sync.WaitGroup, msg chan message, i int) {
			SendMessage(msg, "shine", i)
			fmt.Println("sleep 3 seconds")
			time.Sleep(time.Second * 3)
			defer func() {
				wg.Done()
			}()
		}(wg, msg, i)
	}

	close(msg)

	wg.Add(1)
	go func(wg sync.WaitGroup, msg chan message) {
		defer func() {
			wg.Done()
		}()
		for i := 0; i < 10; i++ {
			for tmp := range msg {
				//tmp.ShowMessage()
				fmt.Println("Now:", time.Now())
				tmp.ShowId()
				//fmt.Println("msg type ", reflect.TypeOf(tmp))
			}
		}
	}(wg, msg)

	wg.Wait()
}
