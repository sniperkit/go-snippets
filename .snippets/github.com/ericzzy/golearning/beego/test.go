package main

import (
    "github.com/astaxie/beego"
)

type MainController struct {
    beego.Controller
}

func (this *MainController) Get() {
    this.Ctx.WriteString(`{"message": "hello", "status": 200}`)
}

func main() {
    //beego.SetLevel(beego.LevelError)
    //log := NewLogger(10000)
    //log.SetLogger("file", `{"filename":"/var/log/info.log"}`)
    beego.SetLogger("file", `{"filename":"/var/log/info.log"}`)
    beego.BeeLogger.DelLogger("console")
    beego.Router("/", &MainController{})
    beego.Run(":8090")
}
