package main

import (
	"github.com/astaxie/beego"
)

type RootController struct {
	beego.Controller
}

func (root *RootController) Get() {
	root.Ctx.WriteString("Appname: " + beego.AppConfig.String("appname") +
		"\nhttpport: " + beego.AppConfig.String("httpport") +
		"\nrunmode: " + beego.AppConfig.String("runmode"))
}

func main() {
	beego.Router("/", &RootController{})
	beego.Run()
}
