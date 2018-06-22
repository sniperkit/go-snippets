package handler

import (
	"net/http"

	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"

	deleteproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/delete"
)

func EventDelete(c *gin.Context) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil || userID == "" {
		httpError(c, http.StatusForbidden, "Could not load user")
		return
	}

	sreq := client.NewRequest("gevvent-event-service", "Delete.Call", &deleteproto.Request{
		EventID: c.Params.ByName("id"),
		UserID:  userID,
	})
	srsp := &deleteproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/")
}
