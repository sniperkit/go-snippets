package main

import (
	server "github.com/asim/go-micro/server"
	"github.com/asim/go-micro/store"
	log "github.com/cihub/seelog"
	_ "github.com/flosch/pongo2-addons"
	"github.com/gin-gonic/gin"
	"github.com/ngerakines/ginpongo2"

	_ "github.com/gophergala/Gevvent/services/gevvent-lib/monitoring"
	"github.com/gophergala/Gevvent/services/gevvent-web/handler"
)

func main() {
	server.Name = "gevvent-web"

	// monitoring.Init(server.Name)
	server.Init()

	// Start web server
	r := gin.New()
	r.Static("/static", "static/")
	r.Use(gin.Recovery())
	r.Use(handler.CheckAuth())
	r.Use(handler.Bugsnag())

	view := r.Group("/")
	view.Use(ginpongo2.Pongo2())
	{
		view.GET("/", handler.Home)
		view.POST("/", handler.Home)

		view.GET("/login", handler.Login)
		view.POST("/login", handler.PostLogin)
		view.GET("/logout", handler.Logout)
		view.GET("/register", handler.Register)
		view.POST("/register", handler.PostRegister)

		view.GET("/events/search", handler.EventSearch)
		view.POST("/events/search", handler.EventSearch)
		view.GET("/events/upcoming", handler.EventUpcomingList)
		view.GET("/events/invitations", handler.EventInvitationsList)
		view.GET("/events/hosting", handler.EventHostingList)
		view.GET("/events/past", handler.EventPastList)
		view.GET("/event/:id", handler.EventView)
		view.GET("/event/:id/invite", handler.EventInvite)
		view.POST("/event/:id/invite", handler.PostEventInvite)
		view.GET("/event/:id/rsvp", handler.EventRSVP)
		view.GET("/event/:id/delete", handler.EventDelete)

		authorised := view.Group("/")
		authorised.Use(handler.Authorised())
		authorised.GET("/events/create", handler.EventCreate)
		authorised.POST("/events/create", handler.PostEventCreate)
	}
	api := r.Group("/api")
	{
		v1 := api.Group("/v1")
		v1.GET("/login", handler.Login)
		v1.GET("/logout", handler.Logout)
		v1.POST("/register", handler.Register)
	}

	r.GET("/ping", handler.Ping)

	r.LoadHTMLTemplates("templats/")

	r.Run(":8080")
}

func loadConfigStr(key string) string {
	item, err := store.Get(key)
	if err != nil {
		log.Warnf("Error loading config value %s from store, %s", key, err)
		return ""
	}

	log.Debugf("Loaded config value %s from store", key)
	return string(item.Value())
}
