package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
)

type PublicController struct {
	beego.Controller
}

func (this *PublicController) Get() {
	v := this.GetSession("hackman")

	if v == nil {
		this.Redirect("/", 302)
		return
	} else {
		w, _ := v.(map[string]string)
		this.Data["Name"] = w["name"]
		this.Data["Avatar"] = w["avatar"]

		hackathons := models.GetAllHackathon()
		this.Data["Hackathons"] = hackathons

		announcements := models.GetAllAnnouncement()
		this.Data["Announcements"] = announcements

		this.TplNames = "public.tpl"
	}
}
