package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	service := ":9009"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}

func checkError(err interface{}) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err)
		os.Exit(1)
	}
}

// 通过上面的代码我们可以看出首先我们根据用户的输入通过net.ResolveTCPAddr获取了一个tcpaddr,
// 然后DialTCP获取了一个TCP链接，然后发送请求信息，最后通过ioutil.ReadAll读取全部的服务器反馈信息。
