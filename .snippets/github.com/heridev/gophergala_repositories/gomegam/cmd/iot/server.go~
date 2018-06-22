package main

import (
	"github.com/astaxie/beego"
	"github.com/gophergala/gomegam/routers/things"
	"log"
	"strconv"
	"time"
)


func RunWeb() {
	log.Printf("IOT starting at %s", time.Now())
	handlerWeb()
	log.Println("IOT killed |_|.")
}



func handlerWeb() {

	beego.SessionOn = true
	beego.SetStaticPath("/static_source", "static_source")
	beego.DirectoryIndex = true
	iot := new(things.ThingsRouter)
	beego.Router("/devices/:id", iot, "get:Get")
	
	//port, _ := config.GetString("beego:http_port")
	//port := utils.GetPort()
	//if port == "" {
		port := "8079"
	//}
	http_port, _ := strconv.Atoi(port)
	beego.HttpPort = http_port
	beego.Run()

}
