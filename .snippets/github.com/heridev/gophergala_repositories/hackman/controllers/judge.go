package controllers

import (
	"github.com/astaxie/beego"
)

// oprations for Judge
type JudgeController struct {
	beego.Controller
}

func (c *JudgeController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("Get", c.Get)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

func (c *JudgeController) Post() {

}

func (c *JudgeController) Get() {

}

func (c *JudgeController) GetAll() {

}

func (c *JudgeController) Put() {

}

func (c *JudgeController) Delete() {

}
