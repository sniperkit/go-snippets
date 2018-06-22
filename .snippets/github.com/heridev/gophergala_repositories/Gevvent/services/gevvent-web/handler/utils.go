package handler

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/cihub/seelog"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"

	"github.com/gophergala/Gevvent/services/gevvent-lib/sessionstore"

	readuserproto "github.com/gophergala/Gevvent/services/gevvent-user-service/proto/readuser"
)

func init() {
	pongo2.RegisterFilter("timestamp", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		return pongo2.AsValue(time.Unix(int64(in.Integer()), 0)), nil
	})
}

func render(c *gin.Context, template string, data map[string]interface{}) {
	sstore, err := sessionstore.New()
	if err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Get user ID
	userID, err := getContextString(c, "userID")
	if err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Load user details
	var user *readuserproto.Response
	if userID != "" {
		user, err = readUser(userID)
		if err != nil {
			httpError(c, http.StatusInternalServerError, err.Error())
			return
		}
		if user != nil {
			data["user_id"] = user.ID
			data["user_name"] = user.Username
		}
	}

	// Load flashes
	var flashes []interface{}
	session, err := sstore.Get(c.Request, "session")
	if err == nil {
		flashes = session.Flashes()
	}

	data["flashes"] = flashes
	data["authenticated"] = user != nil
	data["num_upcoming"] = 0
	data["num_invitations"] = 0

	c.Set("template", template)
	c.Set("data", data)
}

func addFlash(c *gin.Context, flash interface{}) {
	sstore, err := sessionstore.New()
	if err != nil {
		log.Error(err.Error())
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}
	session, err := sstore.Get(c.Request, "session")
	if err != nil {
		log.Error(err.Error())
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	session.AddFlash(flash)
	session.Save(c.Request, c.Writer)
}

func httpError(c *gin.Context, code int, message string) {
	log.Error(message)

	c.Error(fmt.Errorf(message), nil)
	c.String(code, message)
	return
}

func getContextString(c *gin.Context, key string) (string, error) {
	val, err := c.Get(key)
	if err != nil {
		return "", nil
	}

	str, ok := val.(string)
	if !ok {
		return "", nil
	}

	return str, nil
}
