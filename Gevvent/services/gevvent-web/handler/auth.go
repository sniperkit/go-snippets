package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/sessions"

	"github.com/gophergala/Gevvent/services/gevvent-lib/sessionstore"
)

type LoginForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Login(c *gin.Context) {
	errMsg, _ := c.Get("error")
	username, _ := c.Get("username")

	render(c, "templates/auth/login.html", map[string]interface{}{
		"error":    errMsg,
		"username": username,
	})
}

func PostLogin(c *gin.Context) {
	var form LoginForm
	c.BindWith(&form, binding.Form)
	c.Set("username", form.Username)

	if form.Username == "" || form.Password == "" {
		c.Set("error", "You must enter your username and password")
		Login(c)
		return
	}

	token, err := loginUser(form.Username, form.Password)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "Incorrect username or password"):
			c.Set("error", "Incorrect username or password")
			Register(c)
			return
		}
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Save token in session
	sstore, err := sessionstore.New()
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session, err := sstore.Get(c.Request, "auth_session")
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session.Values["token"] = token
	sessions.Save(c.Request, c.Writer)

	redirectURL := "/"
	if v := c.Request.URL.Query().Get("redirect"); v != "" {
		v, err = url.QueryUnescape(v)
		if err == nil {
			redirectURL = v
		}
	}

	c.Redirect(http.StatusFound, redirectURL)
}

func Logout(c *gin.Context) {
	sstore, err := sessionstore.New()
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session, err := sstore.Get(c.Request, "auth_session")
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Invalidate token
	if token, ok := session.Values["token"].(string); ok {
		err := logoutUser(token)
		if err != nil {
			c.Error(err, nil)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	}

	// Clear session
	delete(session.Values, "token")
	sessions.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

type RegisterForm struct {
	Username string `form:"username"`
	Password string `form:"password"`
}

func Register(c *gin.Context) {
	errMsg, _ := c.Get("error")
	username, _ := c.Get("username")

	render(c, "templates/auth/register.html", map[string]interface{}{
		"error":    errMsg,
		"username": username,
	})
}

func PostRegister(c *gin.Context) {
	var form RegisterForm
	c.BindWith(&form, binding.Form)
	c.Set("username", form.Username)
	c.Set("username_invalid", form.Username == "")
	c.Set("password_invalid", form.Username == "")

	if form.Username == "" || form.Password == "" {
		c.Set("error", "You must enter a username and password")
		Register(c)
		return
	}

	_, err := registerUser(form.Username, form.Password)
	if err != nil {
		if strings.Contains(err.Error(), "User already exists") {
			c.Set("error", "User already exists")
			Register(c)
			return
		}

		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := loginUser(form.Username, form.Password)
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Save token in session
	sstore, err := sessionstore.New()
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session, err := sstore.Get(c.Request, "auth_session")
	if err != nil {
		c.Error(err, nil)
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	session.Values["token"] = token
	sessions.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

// CheckAuth checks if a user is logged in. If a user is logged in then their token
// and user ID is added to the request context.
func CheckAuth() gin.HandlerFunc {
	// Save token in session
	sstore, err := sessionstore.New()
	if err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		session, err := sstore.Get(c.Request, "auth_session")
		if err != nil {
			c.Error(err, nil)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		if token, ok := session.Values["token"]; ok && !session.IsNew {
			userID, err := checkAuthorised(token.(string))
			if err == nil {
				c.Set("token", token)
				c.Set("userID", userID)
			} else {
				// Token is not valid so clear
				delete(session.Values, "token")
				sessions.Save(c.Request, c.Writer)
			}
		}

		c.Next()
	}
}

// Authorised ensures that a user is logged in
func Authorised() gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, err := c.Get("userID"); err != nil {
			c.Redirect(http.StatusFound, fmt.Sprintf("/login?redirect=%s", url.QueryEscape(c.Request.URL.Path)))
			return
		}

		c.Next()
	}
}
