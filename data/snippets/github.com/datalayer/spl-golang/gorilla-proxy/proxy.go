package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/datalayer/go-school/gorilla-proxy/handler"
	"github.com/datalayer/go-school/gorilla-proxy/model"
	"github.com/gorilla/mux"
	"github.com/koding/websocketproxy"
)

func main() {

	// Functionalities written in Go.
	/*
		http.HandleFunc("/go", func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprint(w, "This is go.")
		})
	*/
	com := &model.Company{
		Name:  "name-1",
		Tel:   "tel-1",
		Email: "email-1"}
	fmt.Println("%+v", com)
	comj, _ := json.Marshal(com)
	fmt.Println(comj)

	router := mux.NewRouter().StrictSlash(true)
	sub := router.PathPrefix("/go/api/v1").Subrouter()
	sub.Methods("GET").Path("/companies").HandlerFunc(handler.GetCompanies)
	sub.Methods("POST").Path("/companies").HandlerFunc(handler.SaveCompany)
	sub.Methods("GET").Path("/companies/{name}").HandlerFunc(handler.GetCompany)
	sub.Methods("PUT").Path("/companies/{name}").HandlerFunc(handler.UpdateCompany)
	sub.Methods("DELETE").Path("/companies/{name}").HandlerFunc(handler.DeleteCompany)
	http.Handle("/go/", sub)

	// Anything we don't do in Go, we pass to the other platform.
	u, _ := url.Parse("http://localhost:8091")
	http.Handle("/", httputil.NewSingleHostReverseProxy(u))

	w, _ := url.Parse("ws://localhost:8091")
	http.Handle("/ws", websocketproxy.NewProxy(w))

	// Start the server
	http.ListenAndServe(":8080", nil)

}
