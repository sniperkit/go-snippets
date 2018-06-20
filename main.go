package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/search", handleSearch)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func handleSearch(w http.ResponseWriter, r *http.Request) {
	log.Println("serving", r.URL)

	query := r.FormValue("q")
	if query == "" {
		http.Error(w, "Missing 'q' URL parameter", http.StatusBadRequest)
		return
	}
}
