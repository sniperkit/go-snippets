package main

import (
	"net/http"
	"path"
)

func main() {
	http.HandleFunc("/", foo)
	http.ListenAndServe(":3000", nil)
}

func foo(w http.ResponseWriter, r *http.Request) {
	// Assuming you want to serve a photo at 'images/foo.png'
	fp := path.Join("images", "foo.png")
	http.ServeFile(w, r, fp)
}