package main

import (
	"./controller"
	"./middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func main() {
	r := mux.NewRouter()
	dir, _ := os.Getwd()
	fileServer := http.FileServer(http.Dir(filepath.Join(dir, "4-pages", "statics")))
	h := http.StripPrefix("/static/", fileServer)
	r.PathPrefix("/static/").Handler(h)
	r.Handle("/home", &controller.Home{})
	r.Handle("/index", &controller.Home{})
	r.Handle("/", &controller.Home{})

	n := negroni.Classic()
	n.UseHandler(r)
	n.UseFunc(middleware.Route)

	s := http.Server{}
	s.Handler = n
	s.ReadTimeout = time.Second
	s.WriteTimeout = time.Second
	s.Addr = ":8081"
	log.Fatalln(s.ListenAndServe())
}
