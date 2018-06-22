package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Data["ClientID"] = beego.AppConfig.String("client_id")

	v := this.GetSession("hackman")
	if v == nil {
		this.TplNames = "index.tpl"
	} else {
		w, _ := v.(map[string]string)

		if w["profile"] == "admin" {
			this.Redirect("/admin", 302)
			return
		}

		this.Data["Name"] = w["name"]
		this.Data["Avatar"] = w["avatar"]

		hackathons := models.GetAllHackathon()
		this.Data["Hackathons"] = hackathons

		announcements := models.GetAllAnnouncement()
		this.Data["Announcements"] = announcements

		this.TplNames = "public.tpl"
	}
}

func (this *MainController) Logout() {
	this.DelSession("hackman")
	this.Redirect("/", 302)
	return
}
