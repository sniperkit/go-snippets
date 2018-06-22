package main

import (
	"net/http"

	"github.com/gophergala/serradacapivara/handlers"
	"github.com/zenazn/goji"
)

func main() {
	// Website
	goji.Get("/", handlers.Index)
	goji.Get("/search", handlers.Search)
	goji.Get("/map", handlers.Map)
	goji.Get("/site/:id", handlers.Site)
	goji.Get("/about", handlers.About)

	// Static
	goji.Get("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	goji.Serve()
}
