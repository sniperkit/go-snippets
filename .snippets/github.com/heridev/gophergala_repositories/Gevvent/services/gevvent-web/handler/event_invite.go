package handler

import (
	"net/http"
	"strings"

	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	inviteproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/invite"
)

type InviteForm struct {
	Username string `form:"username" binding:"required"`
}

func EventInvite(c *gin.Context) {
	var form InviteForm
	submitted, _ := c.Get("submitted")
	errorMsg, _ := c.Get("error")
	val, err := c.Get("form")
	if err == nil {
		if val, ok := val.(InviteForm); ok {
			form = val
		}
	}

	render(c, "templates/event/invite.html", map[string]interface{}{
		"submitted": submitted,
		"error":     errorMsg,
		"username":  form.Username,
	})
}

func PostEventInvite(c *gin.Context) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil || userID == "" {
		httpError(c, http.StatusForbidden, "Could not load user")
		return
	}

	var form InviteForm
	valid := c.BindWith(&form, binding.Form)
	c.Set("form", form)
	c.Set("submitted", true)

	if !valid {
		c.Set("error", "You must specify a username")
		EventInvite(c)
		return
	}

	sreq := client.NewRequest("gevvent-event-service", "Invite.Call", &inviteproto.Request{
		EventID:     c.Params.ByName("id"),
		UserID:      userID,
		InvitedUser: form.Username,
	})
	srsp := &inviteproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		if strings.Contains(err.Error(), "Could not find user") {
			c.Set("error", "Could not find user")
			EventInvite(c)
			return
		} else {
			httpError(c, http.StatusInternalServerError, err.Error())
			return
		}
	}

	c.Redirect(http.StatusFound, "/event/"+c.Params.ByName("id"))
}
