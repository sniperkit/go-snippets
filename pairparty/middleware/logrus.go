package middleware

import (
	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"time"
)

func Logrus() gin.HandlerFunc {
	var log = logrus.New()
	log.Level = logrus.InfoLevel
	log.Formatter = &logrus.TextFormatter{}
	return func(c *gin.Context) {
		t := time.Now()
		log.WithFields(logrus.Fields{
			"method":  c.Request.Method,
			"request": c.Request.URL,
			"remote":  c.Request.RemoteAddr,
		}).Info("started handling request")
		// before request

		c.Next()

		// after request
		latency := time.Since(t)
		log.WithFields(logrus.Fields{
			"status": c.Writer.Status(),
			"proto":  c.Request.Proto,
			"took":   latency,
		}).Info("completed handling request")
	}
}
