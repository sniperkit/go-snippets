package main

import (
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func main() {

	// goji have DefaultMux and you can create new Mux
	// the goji mux doc: https://godoc.org/github.com/zenazn/goji/web#Mux

	m := web.New()
	goji.Handle("/foo/*", m)

	// http://127.0.0.1:8000/foo/
	// the result  => get2_1.png
	m.Get("/foo/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo!"))
	})

	// http://127.0.0.1:8000/foo/bar
	// the result  => get2_2.png
	m.Get("/foo/bar", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("foo bar!"))
	})

	goji.Serve()
}
