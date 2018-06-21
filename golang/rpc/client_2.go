package main

import (
	"fmt"
	"net/rpc"

	"./common"
)

func main() {
	var args = common.Args{17, 8}
	var result = common.Result{}

	var client, err = rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		fmt.Println("连接RPC服务失败：", err)
	}
	err = client.Call("MathService.Divide", args, &result)
	if err != nil {
		fmt.Println("调用失败：", err)
	}
	fmt.Println("调用结果：", result.Value)
}
