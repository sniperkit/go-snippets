package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
)

type AdminController struct {
	beego.Controller
}

func (this *AdminController) Get() {
	v := this.GetSession("hackman")
	if v == nil {
		this.Redirect("/", 302)
		return
	} else {
		w, _ := v.(map[string]string)
		this.Data["Name"] = w["name"]
		this.Data["Avatar"] = w["avatar"]

		if w["profile"] == "admin" {
			hackathons := models.GetAllHackathon()
			this.Data["Hackathons"] = hackathons

			this.TplNames = "admin.tpl"
		} else if w["profile"] == "user" {
			this.Redirect("/", 302)
			return
		}
	}
}

func (this *AdminController) AdminEvent() {
	v := this.GetSession("hackman")

	if v == nil {
		this.Redirect("/", 302)
		return
	} else {
		w, _ := v.(map[string]string)
		this.Data["Name"] = w["name"]
		this.Data["Avatar"] = w["avatar"]

		if w["profile"] == "admin" {
			hackathon := this.Ctx.Input.Param(":hackathon")
                        this.Data["hackathon"] = hackathon

			this.TplNames = "control.tpl"
		} else if w["profile"] == "user" {
			this.Redirect("/", 302)
			return
		}
	}
}
