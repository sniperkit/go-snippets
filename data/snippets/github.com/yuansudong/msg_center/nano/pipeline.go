package nano

import (
	"fmt"
	"runtime"

	"github.com/go-mangos/mangos/transport/inproc"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/pull"
	"github.com/go-mangos/mangos/protocol/push"
)

// WorkFunc 为工作函数
type WorkFunc func(Args ...interface{}) ([]interface{}, error)

// PullPush 推拉模式
type PullPush struct {
	pull mangos.Socket
	push mangos.Socket
}

// Create 用于创建一个推拉模式
func Create(url string, Func WorkFunc) *PullPush {
	var err error
	PPP := &PullPush{}
	if PPP.pull, err = pull.NewSocket(); err != nil {
		fmt.Println("创建pull 发生了错误 %s", err.Error())
		return nil
	}
	PPP.pull.AddTransport(inproc.NewTransport())
	PPP.pull.AddTransport(ipc.NewTransport())
	PPP.pull.AddTransport(tcp.NewTransport())

	if PPP.push, err = push.NewSocket(); err != nil {
		fmt.Println("创建Push 发生了错误 %s", err.Error())
		return nil
	}
	PPP.pull.AddTransport(inproc.NewTransport())
	PPP.push.AddTransport(ipc.NewTransport())
	PPP.push.AddTransport(tcp.NewTransport())
	if err = PPP.push.Listen(url); err != nil {
		fmt.Println("push监听套接字失败,%s", err.Error())
		return nil
	}
	Num := runtime.NumCPU() * 8
	for i := 0; i < Num; i++ {
		go worker(i, url, PPP.pull, Func)
	}

	return PPP
}

// AddTask 用于一个任务
func (PP *PullPush) AddTask(Args []byte) error {
	return PP.push.Send(Args)
}

func worker(ID int, url string, pull mangos.Socket, Func WorkFunc) {
	var Err error
	var Msg *mangos.Message
	if Err = pull.Dial(url); Err != nil {
		fmt.Printf("%d 连接 %s 失败 ==> %s\n", ID, url, Err.Error())
		return
	}
	for {
		if Msg, Err = pull.RecvMsg(); Err != nil {
			fmt.Printf("%d 发生了错误 %s 失败 ==> %s\n", ID, url, Err.Error())
		}
		Func(Msg.Body)
	}
}
