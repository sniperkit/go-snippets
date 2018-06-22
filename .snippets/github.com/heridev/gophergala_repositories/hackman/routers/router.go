package routers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/controllers"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/logout", &controllers.MainController{}, "get:Logout")
	beego.Router("/callback", &controllers.OauthController{}, "get:ParseCode")
	beego.Router("/admin", &controllers.AdminController{})
	beego.Router("/admin/hackathon/?:hackathon", &controllers.AdminController{}, "get:AdminEvent")
	beego.Router("/organize", &controllers.OrganizeController{}, "post:Organize")

	beego.Router("/announce/?:category", &controllers.AnnounceController{}, "post:Announce")
	beego.Router("/public", &controllers.PublicController{})

	beego.Router("/team", &controllers.TeamController{})
	beego.Router("/confirmteam", &controllers.TeamConfirmController{})
	beego.Router("/declineteam", &controllers.TeamDeclineController{})

	beego.Router("/hackathon", &controllers.HackathonController{})

	beego.Router("/activity/?:hid", &controllers.ActivityController{}, "post:AddActivity")
	beego.Router("/judge", &controllers.JudgeController{})
	beego.Router("/gtcreateteam", &controllers.GithubTeamFormationController{})
}
