package controllers

import (
	"github.com/astaxie/beego"
)

type ActivityController struct {
	beego.Controller
}

func (this *ActivityController) AddActivity() {
  beego.Info(this.Ctx.Input.RequestBody)
  //beego.Error(err)

  this.Data["json"] = "sample"
  this.ServeJson()
}
