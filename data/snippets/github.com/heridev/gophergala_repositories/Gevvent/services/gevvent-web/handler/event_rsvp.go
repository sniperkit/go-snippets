package handler

import (
	"net/http"

	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"

	rsvpproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/rsvp"
)

func EventRSVP(c *gin.Context) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil || userID == "" {
		httpError(c, http.StatusForbidden, "Could not load user")
		return
	}

	answer := rsvpproto.Status_NOT_GOING
	switch c.Request.URL.Query().Get("answer") {
	case "not_going":
		answer = rsvpproto.Status_NOT_GOING
	case "going":
		answer = rsvpproto.Status_GOING
	}

	sreq := client.NewRequest("gevvent-event-service", "RSVP.Call", &rsvpproto.Request{
		EventID: c.Params.ByName("id"),
		UserID:  userID,
		Answer:  answer,
	})
	srsp := &rsvpproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	switch c.Request.URL.Query().Get("answer") {
	case "going":
		c.Redirect(http.StatusFound, "/event/"+c.Params.ByName("id"))
	default:
		c.Redirect(http.StatusFound, "/")
	}
}
