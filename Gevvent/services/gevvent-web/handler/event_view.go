package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	readuserproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/readuser"
)

func EventView(c *gin.Context) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil {
		httpError(c, http.StatusForbidden, "Could not load user")
		return
	}

	// Fetch event
	eventRsp, err := readEvent(c.Params.ByName("id"))
	if err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}
	if eventRsp.Event == nil {
		httpError(c, http.StatusNotFound, "Event not found")
		return
	}

	// Fetch user event data
	var userEventStatus readuserproto.Status
	if userID != "" {
		userEventRsp, err := readUserEvent(c.Params.ByName("id"), userID)
		if err != nil {
			httpError(c, http.StatusInternalServerError, err.Error())
			return
		}

		userEventStatus = userEventRsp.Status
	}
	attendeesRsp, err := getAttendees(c.Params.ByName("id"))
	if err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}
	hostRsp, err := readUser(eventRsp.Event.UserID)
	if err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if eventRsp.Event.GetPrivate() {
		if eventRsp.Event.UserID != userID &&
			userEventStatus != readuserproto.Status_INVITED &&
			userEventStatus != readuserproto.Status_GOING {
			httpError(c, http.StatusForbidden, "You do not have permission to view this event")
			return
		}
	}

	render(c, "templates/event/view.html", map[string]interface{}{
		"event":           eventRsp.Event,
		"host":            hostRsp,
		"attendees":       attendeesRsp.Users,
		"userEventStatus": userEventStatus.String(),
	})
}
