package controllers

import (
	"github.com/astaxie/beego"
	"github.com/pravj/hackman/models"
)

// oprations for Hackathon
type HackathonController struct {
	beego.Controller
}

func (c *HackathonController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// @Title Post
// @Description create Hackathon
// @Param	body		body 	models.Hackathon	true		"body for Hackathon content"
// @Success 200 {int} models.Hackathon.Id
// @Failure 403 body is empty
// @router / [post]
func (c *HackathonController) Post() {

}

// @Title Get
// @Description get Hackathon by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Hackathon
// @Failure 403 :id is empty
// @router /:id [get]
func (c *HackathonController) Get() {
	v := c.GetSession("hackman")
	if v == nil {
		c.Redirect("/", 302)
		return
	}

	hackathonId := c.Input().Get("hackathonId")
	w, _ := v.(map[string]string)
	c.Data["Name"] = w["name"]
	c.Data["Avatar"] = w["avatar"]

	if hackathonId == "" {
		c.Data["Hackathons"] = models.GetAllHackathon()
		c.TplNames = "current.tpl"
		return
	} else {
		c.TplNames = "hackathon.tpl"
		return
	}
	
}

// @Title Get All
// @Description get Hackathon
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Hackathon
// @Failure 403
// @router / [get]
func (c *HackathonController) GetAll() {

}

// @Title Update
// @Description update the Hackathon
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Hackathon	true		"body for Hackathon content"
// @Success 200 {object} models.Hackathon
// @Failure 403 :id is not int
// @router /:id [put]
func (c *HackathonController) Put() {

}

// @Title Delete
// @Description delete the Hackathon
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *HackathonController) Delete() {

}
