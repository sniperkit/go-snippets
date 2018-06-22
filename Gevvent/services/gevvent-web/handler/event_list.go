package handler

import (
	"net/http"

	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"

	listproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/list"
)

func eventList(c *gin.Context, viewType *listproto.ViewType) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil || userID == "" {
		httpError(c, http.StatusForbidden, "Could not load user")
		return
	}

	sreq := client.NewRequest("gevvent-event-service", "List.Call", &listproto.Request{
		UserID:   userID,
		ViewType: viewType,
	})
	srsp := &listproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	render(c, "templates/event/list.html", map[string]interface{}{
		"viewType": viewType.String(),
		"events":   srsp.Events,
	})
}

func EventUpcomingList(c *gin.Context) {
	eventList(c, listproto.ViewType_UPCOMING.Enum())
}
func EventInvitationsList(c *gin.Context) {
	eventList(c, listproto.ViewType_INVITATIONS.Enum())
}
func EventHostingList(c *gin.Context) {
	eventList(c, listproto.ViewType_HOSTING.Enum())
}
func EventPastList(c *gin.Context) {
	eventList(c, listproto.ViewType_PAST.Enum())
}
