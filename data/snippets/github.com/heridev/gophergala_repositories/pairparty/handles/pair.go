package handles

import (
	"github.com/flosch/pongo2"
	"github.com/gin-gonic/gin"
	"net/http"
)

func PairGetHandler(c *gin.Context) {
	ctx := c.MustGet("template_data").(pongo2.Context)
	ctx = ctx.Update(pongo2.Context{
		"title": "Pair",
	})
	c.HTML(http.StatusOK, "templates/pages/pair.html", ctx)
}
