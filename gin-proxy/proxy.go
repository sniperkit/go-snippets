package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/datalayer/go-school/gorilla-rest/handler"
	"github.com/datalayer/go-school/gorilla-rest/model"
)

func main() {

	// Functionalities written in Go.
	/*
		http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "This is go.")
		})
	*/
	com := &model.Company{
		Name:  "name-1",
		Tel:   "tel-1",
		Email: "email-1"}
	fmt.Printf("%+v", com)
	comj, _ := json.Marshal(com)
	print(comj)

	r := gin.Default()
	
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	func ReverseProxy() gin.HandlerFunc {
		target := "http://localhost:8091"	
		return func(c *gin.Context) {
			director := func(req *http.Request) {
				r := c.Request
				req = r
				req.URL.Scheme = "http"
				req.URL.Host = target
				req.Header["my-header"] = []string{r.Header.Get("my-header")}
				// Golang camelcases headers
				delete(req.Header, "My-Header")
			}
			proxy := &httputil.ReverseProxy{Director: director}
			proxy.ServeHTTP(c.Writer, c.Request)
		}
	}

	r.GET("/", ReverseProxy()
	
	port := os.Getenv("PORT")
	port = "8080"
	r.Run(":" + port)
	
}
