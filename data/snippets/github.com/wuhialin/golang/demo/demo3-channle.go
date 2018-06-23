/*
* @Author: wuhailin
* @Date:   2017-10-06 17:25:27
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-12-11 09:08:33
 */
package main

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
)

//var l sync.Mutex
var wait sync.WaitGroup

type Result struct {
	url     string
	content string
}

func main() {
	// l.Lock()
	// go func() {
	// 	for i := 0; i < 1000; i++ {
	// 		println(i)
	// 	}
	// 	l.Unlock()
	// }()
	// l.Lock()
	// println("end")
	//
	// var once sync.Once
	// for i := 0; i < 10; i++ {
	// 	go func() {
	// 		once.Do(test)
	// 	}()
	// }
	//
	// n := 3
	// c := make(chan int, n)
	// defer close(c)
	// for i := 0; i < n; i++ {
	// 	go func() {
	// 		for j := 0; j < 10; j++ {
	// 			println(j)
	// 		}
	// 		c <- i
	// 	}()
	// }

	// select {
	// case <-c:
	// }
	//
	urls := []string{}
	for i := 0; i < 10000; i++ {
		urls = append(urls, "http://127.0.0.1/?wd="+strconv.Itoa(i))
	}
	multiGet(urls)
	// for _, content := range multiGet(urls) {
	// 	println(content)
	// }
	wait.Wait()
}

func multiGet(urls []string) []string {
	urlLen := len(urls)
	urlChan := make(chan Result, urlLen)
	content := []string{}
	defer close(urlChan)
	for _, url := range urls {
		wait.Add(1)
		requestGet(url, urlChan)
	}

	for i := 0; i < urlLen; i++ {
		println(i)
		result := <-urlChan
		content = append(content, result.content)
	}
	return content
}

func requestGet(url string, c chan Result) {
	result := Result{url: url, content: ""}
	response, err := http.Get(url)
	if nil == err {
		defer response.Body.Close()
		body, _ := ioutil.ReadAll(response.Body)
		result.content = string(body)
	}
	c <- result
	wait.Done()
}

// func test() {
// 	for j := 0; j < 100; j++ {
// 		println(j)
// 	}
// }
