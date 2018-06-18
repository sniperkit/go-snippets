package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	serveWeb()

}

func serveWeb() {
	r := mux.NewRouter()

	//Serves static files in directory public
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("public/"))))

	log.Printf("Listening at :8080 . . . \n I am serving to Web \n")
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)

}
