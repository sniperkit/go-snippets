/**
* @Author: wuhailin
* @Date:   2017-09-27 11:58:37
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-09-29 14:40:16
 */

package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	filename := "E:/data/b_story20171.txt"
	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if nil != err {
		log.Fatal(filename + "文件打开失败！")
	}
	defer file.Close()
	r := io.Reader(file)
	pr := bufio.NewReader(r)
	b, _, _ := pr.ReadLine()
	fmt.Println("hello world")
	fmt.Println(string(b))
}
