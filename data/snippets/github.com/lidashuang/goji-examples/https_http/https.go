package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
)

func main() {

	r := web.New()
	//https://127.0.0.1:8000/r
	r.Get("/r", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", "r")
	})

	go graceful.ListenAndServeTLS(":8000", "cert.pem", "key.pem", r)

	r1 := web.New()
	//  http://127.0.0.1:8001/r1
	r1.Get("/r1", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", "r1")
	})

	graceful.ListenAndServe(":8001", r1)
}
