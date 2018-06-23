package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
	"strconv"
	"time"
	"github.com/pravj/hackman/utils/mailer"
)

type TeamController struct {
	beego.Controller
}

type TeamConfirmController struct {
	beego.Controller
}

type TeamDeclineController struct {
	beego.Controller
}

type messageParameter struct {
	email       string
	hackathonId int
	teamName    string
	Invitation int
	Accepted int
	Deleted int
	username string
}

type TeamDetail struct {
	Name string
	User1 string
	User2 string
	User3 string
	User4 string
	AccByU1 int
	AccByU2 int
	AccByU3 int
	AccByU4 int
	Email1 string
	Email2 string
	Email3 string
	Email4 string
}

func (c *TeamConfirmController) URLMapping() {
	c.Mapping("Get", c.Get)
}

func (c *TeamDeclineController) URLMapping() {
	c.Mapping("Get", c.Get)
}

func (c *TeamController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
	// c.Mapping("GetOne", c.GetOne)
	// c.Mapping("GetAll", c.GetAll)
	// c.Mapping("Put", c.Put)
	// c.Mapping("Delete", c.Delete)
}

func (this *TeamController) Post() {
	v := this.GetSession("hackman")
	if v == nil {
		this.Redirect("/", 302)
		return
	} else {
		w, _ := v.(map[string]string)

		username := w["username"]
		name := this.Input().Get("teamName")
		beego.Info("Team is "+name)
		hackathonId, _ := strconv.Atoi(this.Input().Get("hackathonId")) //Add error handling for this

		/*Create team first, then add boys.*/

		/*Get UserId of Current User*/
		Uid1, _ := models.GetUserByUsername(username)
		team, err := models.GetTeamByName(name)

		/*Check wheather team exists*/
		if err == nil {
			/*Update the entries*/
			email := this.Input().Get("email")
			user, uerr := models.GetUserByEmail(email) //Add error handling for this
			if uerr != nil {
				beego.Error(uerr)
			}
			if user == nil {
				this.Data["status"] = 1
				this.Data["errorMessage"] = "No such user Exists!!"
				this.Redirect("/team?hackathonId="+strconv.Itoa(hackathonId), 302)
				return
			}


			tempTeam := team
			if team.UserId2 == -1 {
				tempTeam.UserId2 = user.Id
			} else if team.UserId3 == -1 {
				tempTeam.UserId3 = user.Id
			} else if team.UserId4 == -1 {
				tempTeam.UserId4 = user.Id
			}
			err1 := models.UpdateTeamById(tempTeam)


			if err1 == nil {
				var mParameter messageParameter
				mParameter.email = user.Email
				mParameter.hackathonId = hackathonId
				mParameter.teamName = tempTeam.Name
				mParameter.username = username
				mParameter.Invitation = 1
				mParameter.Accepted = 0
				SendMail(mParameter)
			}
		} else {
			/*Create New team */
			team := models.Team{
				Name:        name,
				UserId1:     Uid1.Id,
				UserId2:     -1,
				UserId3:     -1,
				UserId4:     -1,
				CreatorId:   Uid1.Id,
				AccByU1:     true,
				AccByU2:     false,
				AccByU3:     false,
				AccByU4:     false,
				HackathonId: hackathonId,
				CreatedAt:   time.Now().Local(),
			}
			models.AddTeam(&team)
		}
		this.Redirect("/team?hackathonId="+strconv.Itoa(hackathonId), 302)
	}
}

func (c *TeamController) Get() {
	v := c.GetSession("hackman")
	if v == nil {
		c.Redirect("/", 302)
	} else {
		w, _ := v.(map[string]string)
		c.Data["Name"] = w["name"]
		c.Data["Avatar"] = w["avatar"]
		username := w["username"]
		beego.Info(username)
		hackathonId, _ := strconv.Atoi(c.Input().Get("hackathonId"))
		user, err := models.GetUserByUsername(username)
		if err != nil {
			beego.Error(err)
		}
		team, index := getTeamOfUser(user.Id, hackathonId)

		//beego.Info(team)
		if team == nil {
			c.Data["team"] = 0
			c.Data["index"] = -1
		} else {
			c.Data["team"] = 1
			c.Data["index"] = index

			var teamDetail TeamDetail
			teamDetail.Name = team.Name
			if team.UserId1 != -1 {
				user1, _ := models.GetUserById(team.UserId1)
				teamDetail.User1 = user1.Name
				teamDetail.Email1 = user1.Email
				if team.AccByU1 == true {
					teamDetail.AccByU1 = 1
				} else {
					teamDetail.AccByU1 = 0
				}
			} else {
				teamDetail.User1 = "undefined"
				teamDetail.Email1 = "undefined"
				teamDetail.AccByU1 = 0
			}

			if team.UserId2 != -1 {
				user2, _ := models.GetUserById(team.UserId2)
				teamDetail.User2 = user2.Name
				teamDetail.Email2 = user2.Email
				if team.AccByU2 == true {
					teamDetail.AccByU2 = 1
				} else {
					teamDetail.AccByU2 = 0
				}
			} else {
				teamDetail.User2 = "undefined"
				teamDetail.Email2 = "undefined"
				teamDetail.AccByU2 = 0
			}

			if team.UserId3 != -1 {
				user3, _ := models.GetUserById(team.UserId3)
				teamDetail.User3 = user3.Name
				teamDetail.Email3 = user3.Email
				if team.AccByU3 == true {
					teamDetail.AccByU3 = 1
				} else {
					teamDetail.AccByU3 = 0
				}
			} else {
				teamDetail.User3 = "undefined"
				teamDetail.Email3 = "undefined"
				teamDetail.AccByU3 = 0
			}

			if team.UserId4 != -1 {
				user4, _ := models.GetUserById(team.UserId4)
				teamDetail.User4 = user4.Name
				teamDetail.Email4 = user4.Email
				if team.AccByU4 == true {
					teamDetail.AccByU4 = 1
				} else {
					teamDetail.AccByU4 = 0
				}
			} else {
				teamDetail.User4 = "undefined"
				teamDetail.Email4 = "undefined"
				teamDetail.AccByU4 = 0
			}

			c.Data["teamDetail"] = teamDetail
		}
		beego.Info(team)
		c.Data["hackathonId"] = strconv.Itoa(hackathonId)
		c.TplNames = "team.tpl"
	}
}

func getTeamOfUser(userId int, hackathonId int) (team *models.Team, index int) {

	var err error
	team, err = models.GetTeamByUserId1(userId, hackathonId)
	if team != nil {
		index = 0
		return team, index
	}

	if err != nil {
		beego.Error(err)
	}

	//beego.Info("1")
	team, _ = models.GetTeamByUserId2(userId, hackathonId)
	if team != nil {
		index = 1
		return team, index
	}

	//beego.Info("2")
	team, _ = models.GetTeamByUserId3(userId, hackathonId)
	if team != nil {
		index = 2
		return team, index
	}

	//beego.Info("3")
	team, _ = models.GetTeamByUserId4(userId, hackathonId)
	if team != nil {
		index = 3
		return team, index
	}
	//beego.Info("4")
	//beego.Info(team)
	index = -1
	return nil, index
}

func (this *TeamConfirmController) Get() {
	v := this.GetSession("hackman")
	if v == nil {
		this.Redirect("/", 302)
	} else {
		w, _ := v.(map[string]string)
		hackathonId, _ := strconv.Atoi(this.Input().Get("hackathonId"))
		teamName := this.Input().Get("teamName")
		email := w["email"]

		/*Add hackathon and team validation*/

		user, _ := models.GetUserByEmail(email)
		team, err := models.GetTeamByName(teamName)
		if err != nil {
			beego.Error(err)
		} else {
			beego.Info("No error")
		}

		if team.HackathonId == hackathonId {
			if team.UserId1 == user.Id {
				team.AccByU1 = true
			} else if team.UserId2 == user.Id {
				team.AccByU2 = true
			} else if team.UserId3 == user.Id {
				team.AccByU3 = true
			} else if team.UserId4 == user.Id {
				team.AccByU4 = true
			}

			err := models.UpdateTeamById(team)
			if err == nil {
				var mParameter messageParameter
				mParameter.email = user.Email
				mParameter.hackathonId = hackathonId
				mParameter.teamName = team.Name
				mParameter.username = w["username"]
				mParameter.Invitation = 0
				mParameter.Accepted = 1
				SendMail(mParameter)
			}
		} else {
			this.Data["Status"] = 0
			this.Data["message"] = "This team is not for this hackathon"
		}
		this.Redirect("/team?hackathonId="+strconv.Itoa(hackathonId), 302)
	}
}

func (this *TeamDeclineController) Get() {
	v := this.GetSession("hackman")
	if v == nil {
		this.Redirect("/", 302)
	} else {
		w, _ := v.(map[string]string)
		hackathonId, _ := strconv.Atoi(this.Input().Get("hackathonId"))
		teamName := this.Input().Get("teamName")

		user, _ := models.GetUserByUsername(w["username"])
		team, err := models.GetTeamByName(teamName)
		if err != nil {
			beego.Error(err)
		} else {
			beego.Info("No error")
		}

		if team.HackathonId == hackathonId {
			if team.UserId1 == user.Id {
				team.UserId1 = -1
				team.AccByU1 = false
			} else if team.UserId2 == user.Id {
				team.UserId2 = -1
				team.AccByU2 = true
			} else if team.UserId3 == user.Id {
				team.UserId3 = -1
				team.AccByU3 = true
			} else if team.UserId4 == user.Id {
				team.UserId4 = -1
				team.AccByU4 = true
			}

			err := models.UpdateTeamById(team)
			if err == nil {
				var mParameter messageParameter
				mParameter.email = user.Email
				mParameter.hackathonId = hackathonId
				mParameter.teamName = team.Name
				mParameter.username = w["username"]
				mParameter.Invitation = 0
				mParameter.Accepted = 0
				mParameter.Deleted = 1
				SendMail(mParameter)
			}
		} else {
			this.Data["Status"] = 0
			this.Data["message"] = "This team is not for this hackathon"
		}
		this.Redirect("/team?hackathonId="+strconv.Itoa(hackathonId), 302)

	}
}

func SendMail(mParameter messageParameter) {
	email := mParameter.email
	hackathonId := mParameter.hackathonId
	teamName := mParameter.teamName
	username := mParameter.username

	hackathon, _ := models.GetHackathonById(hackathonId)
	user, _ := models.GetUserByUsername(username)
	mail := mailer.New(beego.AppConfig.String("mailgundomain"), beego.AppConfig.String("apikey"), beego.AppConfig.String("apikeypublic"))
	message := "message"
	if mParameter.Invitation == 1 {
		message = "You have been added in team " + teamName + " By " + user.Name + " for " + hackathon.Name + ",\n <a href=http://172.25.18.220:8080/confirmteam?hackathonId=" + strconv.Itoa(hackathonId) + "&teamName=" + teamName + ">click here</a> to confirm your Team."
		msg := mailer.Message{Heading: "Invitation Mail", Body: message, Receiver: email}
		mail.SendMessage(msg)
	} else if mParameter.Accepted == 1 {
		message = "Your participation has been confirmed in team " + teamName +" in <strong>"+hackathon.Name+"</strong>.\nThank You!"
		msg := mailer.Message{Heading: "Invitation Mail", Body: message, Receiver: email}
		mail.SendMessage(msg)
	} else if mParameter.Deleted == 1 {
		team, _ := models.GetTeamByName(teamName)
		message = user.Name+" has declined your invitation to join the team "+teamName+" for "+hackathon.Name+ "."
		if team.UserId1 != -1 {
			user_d, _ := models.GetUserById(team.UserId1)
			msg := mailer.Message{Heading: "Mail", Body: message, Receiver: user_d.Email}
			beego.Info("Sent mail to "+user_d.Email)
			mail.SendMessage(msg)
		}
		if team.UserId2 != -1 {
			user_d, _ := models.GetUserById(team.UserId2)
			msg := mailer.Message{Heading: "Mail", Body: message, Receiver: user_d.Email}
			beego.Info("Sent mail to "+user_d.Email)
			mail.SendMessage(msg)
		}
		if team.UserId3 != -1 {
			user_d, _ := models.GetUserById(team.UserId3)
			msg := mailer.Message{Heading: "Mail", Body: message, Receiver: user_d.Email}
			beego.Info("Sent mail to "+user_d.Email)
			mail.SendMessage(msg)
		}
		if team.UserId4 != -1 {
			user_d, _ := models.GetUserById(team.UserId4)
			msg := mailer.Message{Heading: "Mail", Body: message, Receiver: user_d.Email}
			beego.Info("Sent mail to "+user_d.Email)
			mail.SendMessage(msg)
		}
	}

	beego.Info(message, email)
  
}
