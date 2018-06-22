package main

import (
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
)

const ImageSize = 90000

var walls map[string]*Wall = make(map[string]*Wall)

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	// Start with a blank wall
	walls["Default"] = &Wall{Name: "Default"}
	go h.run()

	r := mux.NewRouter()
	r.HandleFunc("/", serveHome)
	r.HandleFunc("/ws/{id}", serveWs)
	r.HandleFunc("/error", serveError)

	// Brace yourself, some RESTful AF actions in here

	// Create a new wall
	r.HandleFunc("/wall", newWallHandler).Methods("POST")
	// Add an image to a wall
	r.HandleFunc("/wall/{id}", uploadHandler).Methods("POST")
	// Look at a wall
	r.HandleFunc("/wall/{id}", wallHandler).Methods("GET")

	http.Handle("/", r)

	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("./js"))))
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("./css"))))

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
