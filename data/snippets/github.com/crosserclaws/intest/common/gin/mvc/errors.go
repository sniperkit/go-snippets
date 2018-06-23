package mvc

import (
	ogin "github.com/crosserclaws/intest/common/gin"
	"github.com/gin-gonic/gin"
)

var NotFoundOutputBody = OutputBodyFunc(func(c *gin.Context) {
	ogin.JsonNoRouteHandler(c)
})
