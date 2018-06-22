package handler

import (
	"github.com/bugsnag/bugsnag-go"
	"github.com/gin-gonic/gin"
)

func Bugsnag() gin.HandlerFunc {
	return func(c *gin.Context) {
		notifier := bugsnag.New(c.Request)
		defer notifier.AutoNotify()

		c.Next()

		if err := c.LastError(); err != nil {
			notifier.Notify(err, c.Request)
		}
	}
}
