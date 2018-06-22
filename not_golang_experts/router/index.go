package router

import (
	"github.com/unrolled/render"
	"net/http"
)

func Index(res http.ResponseWriter, req *http.Request) {
	render := render.New(render.Options{
		Directory:     "views",
		Layout:        "layout",
		Extensions:    []string{".tmpl", ".html"},
		Charset:       "UTF-8",
		IsDevelopment: true,
	})
	render.HTML(res, 200, "index", nil)
}

func Dashboard(res http.ResponseWriter, req *http.Request) {
	render := render.New(render.Options{
		Directory:     "views",
		Layout:        "layout",
		Extensions:    []string{".tmpl", ".html"},
		Charset:       "UTF-8",
		IsDevelopment: true,
	})
	render.HTML(res, 200, "dashboard", nil)
}
