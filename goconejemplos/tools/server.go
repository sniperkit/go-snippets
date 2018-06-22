package main

import (
	"flag"
	"log"
	"net/http"
	"os"
)

const (
	root   = `/`
	source = `public`
)

var listen = flag.String("listen", ":8000", "Listen address.")

func init() {
	flag.Parse()
}

type rootHandler struct {
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		return false
	}
	return true
}

func (self *rootHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case `/`:
		http.ServeFile(w, r, source+`/index.html`)
	case `/index.html`, `/site.css`, `/favicon.ico`, `play.png`:
		http.ServeFile(w, r, source+r.URL.Path)
	default:
		rsc := r.URL.Path
		if fileExists(source + `/` + rsc) {
			w.Header().Add(`Content-Type`, `text/html; charset=utf-8`)
			http.ServeFile(w, r, source+rsc)
		} else {
			http.ServeFile(w, r, source+`/404.html`)
		}
	}
}

func main() {
	var err error
	http.Handle(root, &rootHandler{})

	log.Printf("Serving HTTP at: %s", *listen)

	if err = http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("ListenAndServe: %q", err)
	}

}
