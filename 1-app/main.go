package main

import (
	_ "./common"
	"./controller"
	"./middleware"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

func main() {
	runtime.Gosched()
	runtime.GOMAXPROCS(runtime.NumCPU())
	r := mux.NewRouter()
	curDir, _ := os.Getwd()
	h := http.FileServer(http.Dir(filepath.Join(curDir, "public")))
	sp := http.StripPrefix("/public/", h)
	r.PathPrefix("/public/").Handler(sp)
	r.Handle("/vpn", &controller.VPN{})
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
	s.Addr = ":8080"
	log.Fatalln(s.ListenAndServe())
}
