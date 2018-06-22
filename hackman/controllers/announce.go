package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
	"github.com/pravj/hackman/utils/github"
)

type AnnounceController struct {
	beego.Controller
}

type GithubTeamFormationController struct {
	beego.Controller
}

func (c *GithubTeamFormationController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
}

func (c *GithubTeamFormationController) Post() {

}

func (c *GithubTeamFormationController) Get() {

	v := c.GetSession("hackman")
	if v == nil {
		c.Redirect("/", 302)
		return
	}

	hackathonName := c.Input().Get("name")
	hackathon, _ := models.GetHackathonByName(hackathonName)
	teams, err := models.GetAllTeamByHackathonId(hackathon.Id)
	if err != nil {
		beego.Error(err)
	} else {
		beego.Info(len(teams))
	}

	var org github.Organization
	org.Name = hackathon.Name
	var git_teams []github.Team

	for i := 0; i < len(teams); i++ {
        var git_team github.Team
        var git_users []github.User

        git_team.Name = teams[i].Name
        git_team.Id = teams[i].Id

        if teams[i].UserId1 != -1 {
        	var git_user github.User
        	username, _ := models.GetUserById(teams[i].UserId1)
        	git_user.UserName = username.UserName
        	git_users = append(git_users, git_user)
        }
        if teams[i].UserId2 != -1 {
        	var git_user github.User
        	username, _ := models.GetUserById(teams[i].UserId2)
        	git_user.UserName = username.UserName
        	git_users = append(git_users, git_user)
        }
        if teams[i].UserId1 != -1 {
        	var git_user github.User
        	username, _ := models.GetUserById(teams[i].UserId3)
        	git_user.UserName = username.UserName
        	git_users = append(git_users, git_user)
        }
        if teams[i].UserId1 != -1 {
        	var git_user github.User
        	username, _ := models.GetUserById(teams[i].UserId4)
        	git_user.UserName = username.UserName
        	git_users = append(git_users, git_user)
        }
        git_team.Users = git_users
        git_teams = append(git_teams, git_team)
    }
    org.Teams = git_teams
    token, _ := models.GetAdminToken()

    github.CreateTeams(&org, token.Token)
	c.Redirect("/admin/hackathon/"+hackathonName, 302)
	return
}

func (this *AnnounceController) Announce() {
	announcement := this.GetString("announcement")
	category := this.Ctx.Input.Param(":category")

	this.MakeAnnouncement(category, announcement)

        if category != "" {
	  this.Redirect("/admin/hackathon/" + category, 302)
        } else {
          this.Redirect("/admin", 302)
        }
	return
}

func (this *AnnounceController) MakeAnnouncement(category, announcement string) {
	announce := models.Announcement{Category: category, Announcement: announcement}
	models.AddAnnouncement(&announce)
}
