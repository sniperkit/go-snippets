package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
)

func main() {
	goji.Get("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "helloworld..........")
	})

	// custom 404 handler
	// api: https://godoc.org/github.com/zenazn/goji#NotFound
	goji.NotFound(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "the page not found")
	})

	goji.Serve()
}
