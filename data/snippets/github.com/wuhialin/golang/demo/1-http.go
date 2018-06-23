package main

import (
	"./urls"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"time"
)

func main() {
	r := mux.NewRouter()
	r.Handle("/test", &urls.TestHandle{})
	r.Handle("/", &urls.HomeHandle{})

	n := negroni.New()
	h := negroni.HandlerFunc(urls.Middleware)
	n.Use(h)
	n.UseHandler(r)

	s := &http.Server{}
	s.Addr = ":8080"
	s.Handler = n
	s.ReadTimeout = 1 * time.Second
	s.WriteTimeout = 1 * time.Second
	log.Fatalln(s.ListenAndServe())
}
