package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	var (
		port = flag.String("port", "8030", "The WEB Server port number.")
	)
	flag.Parse()

	r := gin.Default()
	r.Use(cors.Default())
	/*
		r.GET("/", func(c *gin.Context) {
			c.Redirect(http.StatusMovedPermanently, "./kuber")
		})
	*/

	//	r.Group("/kuberapi")
	//	{
	r.GET("/api/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	//	}

	r.StaticFS("/", http.Dir("./static"))
	//	r.Static("/assets", "./assets")
	//	r.StaticFile("/favicon.ico", "./resources/favicon.ico")

	fmt.Println("Starting WEB server on port " + *port)

	r.Run(":" + *port)

}
