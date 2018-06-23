package nano

import (
	"fmt"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/rep"
)

type (
	//
	copoll struct {
	}
	// receiver 用于描述一个接受者
	Receiver struct {
		resp mangos.Socket
	}
	// dealer  用于描述一个任务处理着
	Dealer struct {
	}
)

// 创建本地接受者
func createReceiver(url string) *Receiver {
	var err error
	preceiver := &Receiver{}

	if preceiver.resp, err = rep.NewSocket(); err != nil {
		panic(fmt.Sprintf("创建rep模式的套接字失败,失败原因是:%s\n", err.Error()))
	}
	if err = preceiver.resp.Listen(url); err != nil {
		panic(fmt.Sprintf("监听rep模式的套接字失败,失败原因是:%s\n", err.Error()))
	}
	return preceiver
}

// Run 用于开始运行接受者
func (r *Receiver) Run() {
	var err error
	PMessage * mangos.Message
	for {
		if PMessage, err = r.resp.RecvMsg(); err != nil {
			fmt.Println("当前协程接受消息发生了错误,错误的原因是")
		}
		go worker()
	}

}
