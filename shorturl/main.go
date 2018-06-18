package main

import (
	"github.com/astaxie/beego"
	"github.com/jiangew/hancock/shorturl/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/shorten", &controllers.ShortController{})
	beego.Router("/expand", &controllers.ExpandController{})

	beego.Run()
}
