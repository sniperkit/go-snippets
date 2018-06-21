/*
* @Author: wuhialin
* @Date:   2018-01-16 10:00:11
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-22 14:33:07
 */
package main

import (
	"./model"
	"./router"
	"github.com/henrylee2cn/faygo"
)

func main() {
	router.Router(faygo.New(model.AppName))
	faygo.Run()
}
