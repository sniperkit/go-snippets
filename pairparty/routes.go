package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"github.com/jinzhu/gorm"
	"github.com/logie17/Project-V/config"
	h "github.com/logie17/Project-V/handles"
	m "github.com/logie17/Project-V/middleware"
	"github.com/logie17/Project-V/model"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var store *sessions.CookieStore = sessions.NewCookieStore([]byte("a-secret-string"))

func main() {
	router := gin.New()
	configuration := config.LoadConfig()

	db, err := setupDB(configuration)
	if err != nil {
		log.Println("Uh, can't open the database: %s", err.Error())
	}

	router.Use(m.Logrus())
	router.Use(m.IsMobile()) // this doesnt work yet
	// this is how we can get global template data
	//set := pongo2.NewSet("our web templates") // The idea behind sets is that you can create another set with other globals and configurations for mail templates or other kind of templates
	//set.Globals["global_variable"] = "this is a test"
	// https://github.com/flosch/pongo2/issues/35
	router.HTMLRender = newPongoRender()
	router.Static("/public", "./public")

	router.GET("/", h.IndexHandler)

	router.GET("/login", h.LoginGetHandler(store))
	router.POST("/login", h.LoginPostHandler(store, db))

	router.GET("/signup", h.SignupGetHandler(store))
	router.POST("/signup", h.SignupPostHandler(store))

	router.GET("/pair", m.IsAuthenticated(store), h.PairGetHandler)

	router.Run(fmt.Sprintf("[::]:%s", configuration.Port))
}

func setupDB(configuration *config.Configuration) (*gorm.DB, error) {
	db, err := gorm.Open("sqlite3", configuration.Database)

	if err != nil {
		return &db, err
	}

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(100)
	db.SingularTable(true)

	db.CreateTable(&model.User{})

	return &db, err

}
