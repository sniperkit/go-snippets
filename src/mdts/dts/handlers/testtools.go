package handlers

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func printReqBody(req *http.Request) {
	body := req.Body

	bs, err := ioutil.ReadAll(body)
	if err != nil {
		log.Println(err)
		return
	}
	body.Close()

	log.Printf("Body: %s.", string(bs))
}

func printReqHead(req *http.Request) {
	log.Printf("HEAD: %v.", req.Header)
}

func createFilterForSetRecord(unitCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("Record", unitCode)
		c.Next()
	}
}

// RunServerWrapHandler run server with a handle function
func runServerWrapHandler(hostport, path string, handle gin.HandlerFunc) error {
	router := gin.Default()
	router.POST(path, handle)

	return router.Run(hostport)
}

func runServerWrapHandlerMap(hostport string, handles map[string]gin.HandlerFunc) error {
	router := gin.Default()
	router.Use(func(c *gin.Context) {
		c.Next()
	})
	for path, handle := range handles {
		router.POST(path, handle)
	}

	return router.Run(hostport)
}
