package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
	"time"
	"strings"
)

type OrganizeController struct {
	beego.Controller
}

func (this *OrganizeController) Organize() {
	name := this.GetString("hackathon")
	desc := this.GetString("hackathon-description")
	org := this.GetString("hackathon-organization")
	start := this.Input().Get("start-time")
	end := this.Input().Get("end-time")
	format := "2006-01-02 15:04"
	
	arr := strings.Split(start, "T")
    start_t, err := time.Parse(format, arr[0]+" "+arr[1])

    arr = strings.Split(end, "T")
	end_t, _ := time.Parse(format, arr[0]+" "+arr[1])

	if err != nil {
		beego.Error(err)
	}

	hackathon := models.Hackathon{Name: name, Description: desc, Organization: org, CreatedAt: time.Now().Local(), StartedTime: start_t, EndTime: end_t}
	models.CreateHackathon(&hackathon)

	this.Redirect("/admin/hackathon/" + name, 302)
	return
}
