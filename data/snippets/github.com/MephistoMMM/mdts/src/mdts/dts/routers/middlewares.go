package routers

import (
	"bytes"
	"io/ioutil"
	"log"
	pts "mdts/protocols/dtsproto"

	"github.com/gin-gonic/gin"
)

type logResponseWriter struct {
	gin.ResponseWriter
	buf *bytes.Buffer
}

func (lrw *logResponseWriter) Write(p []byte) (int, error) {
	lrw.buf.Write(p)
	return lrw.ResponseWriter.Write(p)
}

// logReqAndRespBody 记录请求和响应日志
func logReqAndRespBody(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.Printf("Error reading body: %v", err)
		c.AbortWithStatusJSON(200, &pts.CommResp{
			Code:    pts.FAILED,
			Message: "Can't Read Request Body",
		})
		return
	}
	log.Println("--------------------------------------------------")
	log.Println("Remote Address: ", c.Request.RemoteAddr)
	// And now set a new body, which will simulate the same data we read:
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	// simulate response
	res := &logResponseWriter{
		ResponseWriter: c.Writer,
		buf:            &bytes.Buffer{},
	}
	c.Writer = res

	log.Println("Request Body: ", string(body))

	c.Next()

	log.Println("Response Body: ", res.buf.String())
	log.Println("--------------------------------------------------")
}
