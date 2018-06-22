package main

import (
	"encoding/json"
	"net/http"
)

type JSONMessage struct {
	Status  int
	Message string
}

func (j *JSONMessage) WriteOut(w http.ResponseWriter) {
	// Stole this from http://www.alexedwards.net/blog/golang-response-snippets
	js, err := json.Marshal(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
