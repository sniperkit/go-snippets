package main

import (
	"net/http"

	"github.com/gorilla/mux"

	"github.com/lcsontos/uwatch/html"
	"github.com/lcsontos/uwatch/util"
	"github.com/lcsontos/uwatch/webservice"
)

type handlerFuncType func(http.ResponseWriter, *http.Request)

func handleSafely(handlerFunc handlerFuncType) handlerFuncType {
	return func(rw http.ResponseWriter, req *http.Request) {
		// Recover Panic
		defer func() {
			if err := recover(); err != nil {
				util.HandlePanic(err, rw, req)
			}
		}()

		// Call original handler
		handlerFunc(rw, req)
	}
}

func init() {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc(
		"/api/long_video_url/{videoTypeName}/{videoId}",
		handleSafely(webservice.GetLongVideoUrl)).Methods("GET")

	router.HandleFunc(
		"/api/parse_video_url",
		handleSafely(webservice.GetParseVideoUrl)).Methods("POST")

	router.HandleFunc(
		"/{urlId}/{urlPath}",
		handleSafely(html.ProcessTemplate)).Methods("GET")

	http.Handle("/", router)
}
