package main

import (
	"fmt"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/zenazn/goji/web"
)

func hello(c web.C, w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, %s!", c.URLParams["name"])
}

func main() {
	// get1: url parameters
	// visit http://127.0.0.1:8000/hello/world
	// the name is parameter
	// the result => get1_1.png
	goji.Get("/hello/:name", hello)

	// Query String Parameters
	// for example visit: http://127.0.0.1:8000/world?name=helloword
	// the result => get1_2.png
	goji.Get("/world", func(c web.C, w http.ResponseWriter, r *http.Request){
		fmt.Fprintf(w, "hello params is %s", r.FormValue("name"))
	})

	goji.Serve()
}
