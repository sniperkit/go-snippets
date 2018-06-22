package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
	"net/http"
)

func IsAuthenticated(store *sessions.CookieStore) gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "flash-session")
		if err != nil {
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		}
		var email = session.Values["email"]
		if email == nil {
			session.AddFlash("Login failed!", "message")
			session.Save(c.Request, c.Writer)
			c.Fail(http.StatusUnauthorized, errors.New("Unauthorized")) // idk why this is needed but it is
			c.Redirect(http.StatusMovedPermanently, "/login")
		}
	}
}
