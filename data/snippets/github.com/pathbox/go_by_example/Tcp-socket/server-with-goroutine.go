package main

import (
	"fmt"
	"net"
	"os"
	// "time"
)

func main() {
	service := ":9009"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	fmt.Println(tcpAddr)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	for {
		conn, err := listener.Accept()
		fmt.Println(conn)
		if err != nil {
			continue
		}
		go handlerClient(conn)
	}
}

func handlerClient(conn net.Conn) {
	defer conn.Close()
	daytime := "hello kitty"
	conn.Write([]byte(daytime)) // don't care about return value
	// we're finished with this client
}

func checkError(err interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
		os.Exit(1)
	}
}

// 上面的服务我们跑起来之后，他将会一直在那边等待，直到有新的客户端到达，当有新的客户端到达的时候他反馈当前的时间信息。同时我们注意看循环那里，当有错误发生时，直接continue了，
// 而不是退出，因为在我们编写服务器端的时候，当有错误发生的情况下最好是记录错误，然后当前客户端出错直接退出，而不会影响到当前整个的服务。
