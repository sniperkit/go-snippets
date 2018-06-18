package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

// Login binding from JSON
type Login struct {
	User string `form:"user" json:"user" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	// curl -v -X POST \
  	// http://localhost:8080/loginJSON \
  	// -H 'content-type: application/json' \
  	// -d '{ "user": "jiangew" }'
	// Example for binding JSON ({"user": "jiangew", "password": "123"})
	r.POST("/loginJSON", func(c *gin.Context) {
		var json Login
		if err := c.ShouldBindJSON(&json); err == nil {
			if json.User == "jiangew" && json.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "You are logined in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		}
	})

	// Example for binding a HTML form (user=jiangew&password=123)
	r.POST("/loginForm", func(c *gin.Context) {
		var form Login
		// This will infer what binder to use depending on the content-type header.
		if err := c.ShouldBind(&form); err == nil {
			if form.User == "manu" && form.Password == "123" {
				c.JSON(http.StatusOK, gin.H{"status": "you are logged in"})
			} else {
				c.JSON(http.StatusUnauthorized, gin.H{"status": "unauthorized"})
			}
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
	})

	// Redirect
	r.GET("/redirect", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "www.google.com")
	})
}

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
