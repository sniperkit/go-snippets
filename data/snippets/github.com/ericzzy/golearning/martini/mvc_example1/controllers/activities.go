package controllers

import (
	"../models"
	"fmt"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	//"log"
	//"net/http"
)

func IndexActivities(render render.Render, db *gorm.DB) {
	activities := []models.Activity{}
	db.Find(&activities)
	fmt.Printf("%+v\n", activities)
	render.HTML(200, "activities/index", activities)
}

func ShowActivity(params martini.Params, render render.Render, db *gorm.DB) {
	activity := models.Activity{}
	db.First(&activity, params["id"])
	render.HTML(200, "activities/show", activity)
}

func NewActivity(render render.Render) {
	render.HTML(200, "activities/new", models.Activity{})
}

func EditActivity(params martini.Params, render render.Render, db *gorm.DB) {
	activity := models.Activity{}
	db.First(&activity, params["id"])
	render.HTML(200, "activities/edit", activity)
}

func CreateActivity(params martini.Params, render render.Render, activity models.Activity, db *gorm.DB) {
	db.Save(&activity)
	render.Redirect("/activities/"+string(activity.ID), 302)
}

func UpdateActivity(params martini.Params, render render.Render, activity_updated models.Activity, db *gorm.DB) {
	activity := models.Activity{}
	db.First(&activity, params["id"])
	db.Model(&activity).Updates(activity_updated)
	render.Redirect("/activities/"+string(params["id"]), 302)
}

func DeleteActivity(params martini.Params, render render.Render, db *gorm.DB) {
	db.Delete(models.Activity{}, params["id"])
	render.Redirect("/activities/", 302)
}
