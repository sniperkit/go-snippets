package main

import (
	"fmt"
	"otp"
	"time"
)

var debugTest string = "********(test)*********> "

type test struct {
	name string
	*otp.GenServer
}

func (t *test) Init() (string, interface{}, error) {
	fmt.Println(debugTest, "test init")
	return otp.NORMAL, 1000, nil
}

func (t *test) HandleMessage(message interface{}) (string, interface{}, error) {
	fmt.Println(debugTest, "test handle message:", message.(string))
	t.name = message.(string)
	return otp.NOREPLY, message, nil
}

func (t *test) HandleCall(msg interface{}) (string, interface{}, error) {
	fmt.Println(debugTest, "test handle call")
	time.Sleep(time.Second * 3)
	return otp.REPLY, msg, nil
}

func (t *test) HandleInfo(msg interface{}) (string, int, error) {
	fmt.Println(debugTest, "test handle cast")
	return otp.NOREPLY, 2000, nil
}

func (t *test) Stop(reason string) error {
	fmt.Println(debugTest, "test handle stop")
	return nil
}

func (t *test) Terminate(reason string) {
	fmt.Println(debugTest, "test handle terminate:", reason)
	return
}

func main() {
	fmt.Println(debugTest, "start gen_server")

	t := &test{}
	genServer := otp.NewGenServer(t)
	err := genServer.Start()
	if err != nil {
		fmt.Println(debugTest, "test start gen_server error:", err)
		return
	}

	genServer.Cast("hello world")
	fmt.Println(debugTest, "test send cast info end:", time.Now())

	//res, err := genServer.Call("gen_server call", 4000)
	//res, err := genServer.Call("gen_server call", otp.INFINITY)
	fmt.Println(debugTest, "test start call :")
	res, err := genServer.Call("(foo call)", nil)
	if err != nil {
		fmt.Println(debugTest, "test call error:", err)
		return
	}
	fmt.Println(debugTest, "test call result:", res.(string), time.Now())

	stop := make(chan string, 1)
	go func() {
		time.Sleep(time.Second * 10)
		stop <- "stop"
	}()

	select {
	case <-stop:
		fmt.Println(debugTest, "test stop")
	}

	fmt.Println(debugTest, "test name", t.name)
}
