package middleware

import (
	"github.com/Shaked/gomobiledetect"
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
)

func IsMobile() gin.HandlerFunc {
	return func(c *gin.Context) {
		detect := mobiledetect.NewMobileDetect(c.Request, nil)
		// We'll have to define this further upstream
		ctx := pongo2.Context{"mobile": false}
		if detect.IsMobile() {
			ctx.Update(pongo2.Context{"mobile": true})
		}
		c.Set("template_data", ctx)
		c.Next()
	}
}
