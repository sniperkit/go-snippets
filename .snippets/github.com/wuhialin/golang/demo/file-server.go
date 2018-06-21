/*
* @Author: wuhailin
* @Date:   2017-09-29 15:08:12
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-09-29 15:43:42
 */
package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
)

func main() {
	path := flag.String("dir", "", "指定监听路径")
	port := flag.Int("port", 0, "指定监听端口")
	flag.Parse()
	if "" == *path {
		log.Fatal("指定文件路径")
	}
	if 0 == *port {
		log.Fatal("请指定监听端口")
	}
	h := http.FileServer(http.Dir(*path))
	err := http.ListenAndServe(":"+strconv.Itoa(*port), h)
	if nil != err {
		log.Fatal(err)
	}
}
