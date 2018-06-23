package main

import (
	"fmt"
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gophergala/vebomm/core"
	"gopkg.in/validator.v2"
)

func api(r *gin.RouterGroup) {
	r.POST("/register", register)
	r.POST("/login", login)
	r.GET("/checkonline/:id", checkonline)
	r.GET("/rate/:id/:rating", rate)
}

func getId(idStr string) int64 {
	var id int64
	fmt.Sscan(idStr, &id)
	return id
}

func rate(c *gin.Context) {
	id := getId(c.Params.ByName("id"))
	var rating int
	fmt.Sscan(c.Params.ByName("rating"), &rating)
	var user core.User
	err := Db().SelectOne(&user, "SELECT * FROM users WHERE id=$1", id)
	if err != nil {
		fmt.Println(err)
		return
	}
	newMMR := user.Mmr+rating*rating/4
	Db().Exec("UPDATE users SET mmr=$1 WHERE id=$2", newMMR, id)
	
	c.JSON(200, newMMR)
}

func checkonline(c *gin.Context) {
	idStr := c.Params.ByName("id")
	id := getId(idStr)
	if _, ok := h.onlineUsers[id]; ok {
		c.JSON(200, true)
		return
	}
	
	c.JSON(200, false)
}

func register(c *gin.Context) {
	var user core.User
	err := c.ParseBody(&user)
	if err != nil {
		c.Abort(http.StatusBadRequest)
		return
	}

	err = Db().Insert(&user)
	if err != nil {
		c.Abort(http.StatusInternalServerError)
		log.Println(err)
		return
	}

	if validator.Validate(user) != nil {
		c.JSON(200, core.RegisterResult{ValOk: false})
		return
	}

	c.JSON(200, checkDup(user))
}

func checkDup(u core.User) core.RegisterResult {
	var dupUsername int64
	var dupEmail int64
	var err error
	if u.Username != "" {
		dupUsername, err = Db().SelectInt("SELECT COUNT(*) FROM users WHERE username=$1", u.Username)
		if err != nil {
			panic(err.Error())
		}
	}

	if u.Email != "" {
		dupEmail, err = Db().SelectInt("SELECT COUNT(*) FROM users WHERE email=$1", u.Email)
		if err != nil {
			panic(err.Error())
		}
	}

	return core.RegisterResult{
		true,
		dupUsername != 0,
		dupEmail != 0,
	}
}

func login(c *gin.Context) {
	var user core.User
	err := c.ParseBody(&user)
	if err != nil {
		c.Abort(http.StatusBadRequest)
		return
	}

	if user.Id != 0 {
		err = Db().SelectOne(&user, "SELECT * FROM users WHERE id=$1", user.Id)
	} else {
		err = Db().SelectOne(&user, "SELECT * FROM users WHERE email=$1 AND password=$2", user.Email, user.Password)
	}

	if err == sql.ErrNoRows {
		c.JSON(200, core.LoginResult{
			Ok:   false,
			User: nil,
		})
		return
	}

	if err != nil {
		panic(err.Error())
	}

	c.JSON(200, core.LoginResult{
		Ok:   true,
		User: &user,
	})
}
