package handlers

import (
	"io/ioutil"
	pts "mdts/protocols/dtsproto"

	"github.com/gin-gonic/gin"
)

const (
	thirdPath   = "http://127.0.0.1:9000/"
	servicePath = "http://127.0.0.1:9010/"
)

// Echo ...
func Echo(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, "Failed read request body")
		return
	}
	c.JSON(200, &pts.CommResp{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    data,
	})
}
