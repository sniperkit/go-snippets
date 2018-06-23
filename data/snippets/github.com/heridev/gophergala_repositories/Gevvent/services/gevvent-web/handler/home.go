package handler

import (
	"net/http"

	"code.google.com/p/goprotobuf/proto"
	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"

	newestproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/newest"
)

func Home(c *gin.Context) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	if userID == "" {
		homePublic(c)
	} else {
		homeAuthorised(c)
	}
}

func homePublic(c *gin.Context) {
	sreq := client.NewRequest("gevvent-event-service", "Newest.Call", &newestproto.Request{
		Count: proto.Int64(10),
	})
	srsp := &newestproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	render(c, "templates/index.html", map[string]interface{}{
		"events": srsp.Events,
	})
}

func homeAuthorised(c *gin.Context) {
	sreq := client.NewRequest("gevvent-event-service", "Newest.Call", &newestproto.Request{
		Count: proto.Int64(10),
	})
	srsp := &newestproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	render(c, "templates/index.html", map[string]interface{}{
		"events": srsp.Events,
	})
}
