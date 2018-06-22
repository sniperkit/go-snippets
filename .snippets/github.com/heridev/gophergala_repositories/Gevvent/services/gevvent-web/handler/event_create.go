package handler

import (
	"code.google.com/p/goprotobuf/proto"
	"net/http"
	"time"

	"github.com/asim/go-micro/client"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	createproto "github.com/gophergala/Gevvent/services/gevvent-event-service/proto/create"
)

type CreateForm struct {
	Name        string  `form:"name" binding:"required"`
	Description string  `form:"description" binding:"required"`
	WhenFrom    string  `form:"whenfrom" binding:"required"`
	WhenTo      string  `form:"whento" binding:"required"`
	WhereLat    float64 `form:"wherelat" binding:"required"`
	WhereLng    float64 `form:"wherelng" binding:"required"`
	WhereAddr   string  `form:"whereaddr" binding:"required"`
	Private     string  `form:"private"`
	PublicAddr  string  `form:"publicaddr"`
}

func EventCreate(c *gin.Context) {
	var form CreateForm
	submitted, _ := c.Get("submitted")
	errorMsg, _ := c.Get("error")
	val, err := c.Get("form")
	if err == nil {
		if val, ok := val.(CreateForm); ok {
			form = val
		}
	}

	render(c, "templates/event/create.html", map[string]interface{}{
		"submitted":   submitted,
		"error":       errorMsg,
		"name":        form.Name,
		"description": form.Description,
		"whenfrom":    form.WhenFrom,
		"whento":      form.WhenTo,
		"wherelat":    form.WhereLat,
		"wherelng":    form.WhereLng,
		"whereaddr":   form.WhereAddr,
		"private":     form.Private,
		"publicaddr":  form.PublicAddr,
	})
}

func PostEventCreate(c *gin.Context) {
	// Get user ID
	var userID string
	userID, err := getContextString(c, "userID")
	if err != nil || userID == "" {
		httpError(c, http.StatusForbidden, "Could not load user")
		return
	}

	var form CreateForm
	valid := c.BindWith(&form, binding.Form)
	c.Set("form", form)
	c.Set("submitted", true)

	if !valid {
		c.Set("error", "You must fill all the required fields")
		EventCreate(c)
		return
	}

	whenFrom, errFrom := time.Parse("2 January 2006, 15:04:05 -07:00", form.WhenFrom)
	whenTo, errTo := time.Parse("2 January 2006, 15:04:05 -07:00", form.WhenTo)

	if errFrom != nil || errTo != nil || !whenFrom.After(time.Now()) || !whenTo.After(whenFrom) {
		c.Set("error", "The entered time range is invalid")
		EventCreate(c)
		return
	}

	sreq := client.NewRequest("gevvent-event-service", "Create.Call", &createproto.Request{
		UserID:      userID,
		Name:        form.Name,
		Description: form.Description,
		When: createproto.TimeRange{
			From: whenFrom.Unix(),
			To:   whenTo.Unix(),
		},
		Where: createproto.Location{
			Address: form.WhereAddr,
			Lat:     form.WhereLat,
			Lng:     form.WhereLng,
		},
		Private:    proto.Bool(form.Private == "on"),
		PublicAddr: proto.Bool(form.PublicAddr == "on"),
	})
	srsp := &createproto.Response{}

	// Call service
	if err := client.Call(sreq, srsp); err != nil {
		httpError(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.Redirect(http.StatusFound, "/event/"+srsp.ID)
}
