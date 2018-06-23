package main

import (
	"crypto/md5"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
)

var workerNum int = 32
var defaultCount int = 10000
var hosturi string = "http://127.0.0.1:9880/ping"

var defaultClient fasthttp.Client

func init() {
	defaultClient = fasthttp.Client{}
	defaultClient.Name = "defaultClient"
	defaultClient.MaxConnsPerHost = 1024
}

func convertMd516bytesToString(data [16]byte) string {
	var md5_ string
	for _, c := range data {
		part := fmt.Sprintf("%x", c)
		md5_ += part
	}
	return md5_
}

func worker(count int) error {
	for count > 0 {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(hosturi)
		req.Header.SetMethod("POST")

		resp := fasthttp.AcquireResponse()
		if err := defaultClient.DoTimeout(req, resp, time.Second*3); err != nil {
			return err
		}
		body := resp.Body()
		md5sum := convertMd516bytesToString(md5.Sum(body))

		// server md5sum
		md5sumServer := resp.Header.Peek("X-Data-Md5")
		if string(md5sumServer) != md5sum {
			return fmt.Errorf("data md5(%s) is different with server(%s)", md5sum, string(md5sumServer))
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		count--
	}
	return nil
}

func main() {
	var wg sync.WaitGroup
	wg.Add(workerNum)

	for i := 0; i < workerNum; i++ {
		go func() {
			defer wg.Done()
			err := worker(defaultCount)
			if err != nil {
				panic(err)
			}
			fmt.Println("worker done")
		}()
	}

	wg.Wait()
}
