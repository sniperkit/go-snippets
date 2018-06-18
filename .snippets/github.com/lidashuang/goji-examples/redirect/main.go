package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
)

func main() {
	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "root page")
	})

	goji.Get("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "hello page")
	})

	// redirect handler
	// use http.RedirectHandler
	// api: http://godoc.org/net/http#RedirectHandler
	goji.Get("/world", http.RedirectHandler("/hello", 301))

	goji.Serve()
}
