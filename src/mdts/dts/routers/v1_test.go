package routers

import (
	"io/ioutil"
	pts "mdts/protocols/dtsproto"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parnurzeal/gorequest"
)

func echo(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.String(500, "Failed read request body")
		return
	}
	c.JSON(200, &pts.RespT2S{
		Code:    pts.SUCCESS,
		Message: "",
		Data:    string(data),
	})
}

// RunServerWrapHandler run server with a handle function
func runServerWrapMiddle(hostport, path string, middle gin.HandlerFunc) error {
	router := gin.Default()
	router.Use(middle)
	router.POST(path, echo)

	return router.Run(hostport)
}

func TestRouter(t *testing.T) {
	var respBody pts.RespT2S
	req := gorequest.New()

	// success
	reqRouterSuccess := func() {
		resp, _, errs := req.
			Post("http://localhost:8086/testrouter").Type("json").
			Send(`{"Router": 1}`).EndStruct(&respBody)
		if errs != nil {
			t.Error(errs)
			return
		}
		if statuscode := resp.StatusCode; statuscode != 200 {
			t.Errorf("Hope Get Status %d, but get %d.", 200, statuscode)
		}
		if code := respBody.Code; code != pts.SUCCESS {
			t.Errorf("Hope Get Code %d, but get %d.", pts.SUCCESS, code)
			t.Errorf("Message: %s.", respBody.Message)
		}
		if data := respBody.Data; data != `{"Router":1}` {
			t.Errorf("Hope Get Data `{\"Router\":1}`, But Get %s.", data)
		}
	}

	go runServerWrapMiddle(":8086", "/testrouter", logReqAndRespBody)
	time.Sleep(time.Millisecond * 100)
	reqRouterSuccess()
}
