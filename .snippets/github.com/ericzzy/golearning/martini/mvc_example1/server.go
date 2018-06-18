package main

import (
	"log"

	"github.com/codegangsta/martini-contrib/render"
	"github.com/go-martini/martini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/martini-contrib/auth"
	"github.com/martini-contrib/binding"

	"./controllers"
	"./models"
)

func checkErr(err error, message string) {
	if err != nil {
		log.Fatal(message)
	}
}

func main() {
	db, err := gorm.Open("mysql", "root:passw0rd@/martini_mvc?charset=utf8&parseTime=True&loc=Local")
	checkErr(err, "Database connection failed")
	m := martini.Classic()
	m.Map(&db)
	m.Use(render.Renderer(render.Options{
		Layout: "layout",
	}))
	m.Get("/", func() string {
		return "Welcome to ActivityApp"
	})
	m.Group("/activities", func(r martini.Router) {
		r.Get("/", controllers.IndexActivities)
		r.Get("/new", controllers.NewActivity)
		r.Post("/", binding.Bind(models.Activity{}), controllers.CreateActivity)
		r.Get("/:id/edit", controllers.EditActivity)
		r.Post("/:id", controllers.UpdateActivity)
		r.Get("/:id", controllers.ShowActivity)
		r.Get("/:id/delete", controllers.DeleteActivity)
	})
	m.RunOnAddr(":9099")
	m.Run()
}
