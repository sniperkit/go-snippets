package handler

import (
	"net/http"
	"strconv"

	"code.google.com/p/goprotobuf/proto"
	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"

	searchproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/search"
)

func EventSearch(c *gin.Context) {
	// Try to get user ID
	userID, _ := getContextString(c, "userID")

	var location *searchproto.Location
	if c.Request.URL.Query().Get("lat") != "" && c.Request.URL.Query().Get("lng") != "" {
		lat, _ := strconv.ParseFloat(c.Request.URL.Query().Get("lat"), 64)
		lng, _ := strconv.ParseFloat(c.Request.URL.Query().Get("lng"), 64)

		location = &searchproto.Location{
			Lat: lat,
			Lng: lng,
		}
	}

	sreq := client.NewRequest("gevvent-event-service", "Search.Call", &searchproto.Request{
		UserID: proto.String(userID),
		Query:  c.Request.URL.Query().Get("query"),
		Where:  location,
	})
	srsp := &searchproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	render(c, "templates/event/search.html", map[string]interface{}{
		"query":  c.Request.URL.Query().Get("query"),
		"lat":    c.Request.URL.Query().Get("lat"),
		"lng":    c.Request.URL.Query().Get("lng"),
		"addr":   c.Request.URL.Query().Get("addr"),
		"events": srsp.Events,
	})
}
