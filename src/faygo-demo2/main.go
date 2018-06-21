package main

import (
	"github.com/henrylee2cn/faygo"
	"faygo-demo2/models"
	"faygo-demo2/router"
)

func main() {
	frame := faygo.New(models.AppName)
	router.Router(frame)
	faygo.Run()
}
