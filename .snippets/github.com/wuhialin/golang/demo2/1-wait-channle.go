/*
* @Author: wuhailin
* @Date:   2017-12-11 09:13:00
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-12-11 10:59:44
 */

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"
)

const (
	TASK_MAX_QUERY = 1000
	MAX_URL        = 100000
)

var wait sync.WaitGroup
var f *os.File

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	runtime.Gosched()
	var urls []string
	var err error
	f, err = os.OpenFile("1.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rage := 0
	ch := make(chan bool, TASK_MAX_QUERY)
	for i := 0; i < MAX_URL; i++ {
		rage = i % TASK_MAX_QUERY
		urls = append(urls, fmt.Sprintf("http://127.0.0.1/?wd=%d", i))
		if rage == TASK_MAX_QUERY-1 || i == MAX_URL-1 {
			if len(urls) > 0 {
				multiQuery(urls, ch)
				urls = urls[:0]
			}
		}
		println(i)
	}
	wait.Wait()
}

func multiQuery(urls []string, ch chan bool) {
	for _, url := range urls {
		ch <- true
		wait.Add(1)
		go query(url, ch)
	}
}

func query(url string, ch <-chan bool) {
	var html string
	defer func() {
		if html != "" {
			f.WriteString(fmt.Sprintln(html))
		}
		<-ch
		wait.Done()
	}()
	response, err := http.Get(url)
	if err != nil {
		html = err.Error()
		return
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		html = err.Error()
		return
	}
	html = string(bytes)
	html = strings.Replace(html, "\n", "", -1)
	return
}
