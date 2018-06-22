package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"github.com/gophergala/sleuth/sleuth"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const clientId = "b005064ced134b19b31a9f5317d4fa8c"

var renderer *render.Render

func main() {
	renderer = render.New()
	sleuth.Init(0)
	port := getPort()
	router := mux.NewRouter()
	api := router.PathPrefix("/api/").Subrouter()
	api.HandleFunc("/search", search).Methods("GET").Queries("tags", "{[a-zA-Z0-9_\\s]}")
	api.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "API is running!\n")
	})

	n := negroni.Classic()
	n.UseHandler(api)
	n.Run(":" + port)
}

// Gets the port to bind to when starting the server.
func getPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	return port
}

func search(w http.ResponseWriter, r *http.Request) {
	var result *sleuth.Result
	unesc, _ := url.QueryUnescape(r.URL.RawQuery)
	vars, _ := url.ParseQuery(unesc)
	tags := strings.Split(vars["tags"][0], ",")
	token := vars["token"]
	if len(token) > 0 && token[0] != "" {
		result = sleuth.AuthenticatedSearch(tags, token[0])
	} else {
		result = sleuth.UnauthenticatedSearch(tags, clientId)
	}
	renderer.JSON(w, http.StatusOK, *result)
}
