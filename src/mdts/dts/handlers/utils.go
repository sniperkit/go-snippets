package handlers

import (
	"log"

	"github.com/gin-gonic/gin"
)

func logError(c *gin.Context, err error) {
	log.Println(c.Request.RemoteAddr, c.Request.RequestURI, err)
}
