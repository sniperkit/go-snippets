/**
 * 并发超时处理
 * Go 通过 Select 可以同时处理多个 Channel, Select 默认是阻塞的；只有当监听的 Channel 中有发送或接收可以进行时才会运行，
 * 当同时有多个可用的 Channel, Select 按随机顺序进行处理, Select可以方便处理多 Channel 同时响应；
 * 在 goroutine 阻塞的情况也可以方便借助 Select 超时机制来解除阻塞僵局
 */
package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	url := "jiangew.github.io"
	content, ok := getHTTPRes(url)
	fmt.Println("content: ", content, " status:", ok)
}

// getHTTPRes 获取url的访问值，返回值:1.成功,返回Body部分,2.失败 返回err 3.超时
func getHTTPRes(url string) (string, error) {
	res := make(chan *http.Response, 1)
	httpError := make(chan *error)
	go func() {
		resp, err := http.Get(url)
		if err != nil {
			httpError <- &err
		}
		res <- resp
	}()

	for {
		select {
		case r := <-res:
			result, err := ioutil.ReadAll(r.Body)
			defer r.Body.Close()
			return string(result), err

		case err := <-httpError:
			return "err", *err

		case <-time.After(2000 * time.Millisecond):
			return "Timed out", errors.New("Time out")
		}
	}
}
