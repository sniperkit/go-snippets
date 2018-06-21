package main

import (
	"github.com/henrylee2cn/faygo"
	"faygo-demo1/model"
	"faygo-demo1/router"
	"faygo-demo1/migration"
)

func main() {
	frame := faygo.New(model.AppName)
	router.Router(frame)
	go migration.Run()
	faygo.Run()
}
