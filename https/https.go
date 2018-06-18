package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/zenazn/goji/graceful"
	"github.com/zenazn/goji/web"
	"github.com/zenazn/goji/web/middleware"
)

func main() {
	// goji serve with https doc:
	// http://godoc.org/github.com/zenazn/goji/graceful#ListenAndServeTLS
	// issue: https://github.com/zenazn/goji/issues/40
	// example code:

	r := web.New()
	// use middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", IndexHandler)

	// func ListenAndServeTLS(addr, certfile, keyfile string, handler http.Handler) error
	// Notice, the cert.pem & key.pem is generate for testing
	log.Fatal(graceful.ListenAndServeTLS(":8000", "cert.pem", "key.pem", r))
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", "world")
}
