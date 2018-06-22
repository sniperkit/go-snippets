package handles

import (
	"github.com/Sirupsen/logrus"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/sessions"
	"net/http"
)

type SignupForm struct {
	FirstName       string `form:"f_name" binding:"required"`
	LastName        string `form:"l_name" binding:"required"`
	Organization    string `form:"organization" binding:"required"`
	Email           string `form:"email" binding:"required"`
	Password        string `form:"password" binding:"required"`
	PasswordConfirm string `form:"password_confirm" binding:"required"`
}

func SignupGetHandler(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := pongo2.Context{
			"title": "Pair",
		}
		c.HTML(http.StatusOK, "templates/pages/signup.html", ctx)
	}
}

func SignupPostHandler(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		// lets checkout the form
		var form SignupForm
		c.BindWith(&form, binding.Form)
		logrus.WithFields(logrus.Fields{
			"FirstName":    form.FirstName,
			"LastName":     form.LastName,
			"Organization": form.Organization,
			"Email":        form.Email,
			"Password":     form.Password,
		}).Info("User signup")
		// logie should do something about this
		ctx := pongo2.Context{
			"title": "Pair",
		}
		c.HTML(http.StatusOK, "templates/pages/signup.html", ctx)
	}
}
