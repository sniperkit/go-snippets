package main

import (
	"encoding/json"
	"github.com/JarekSed/stride/parser"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

type request struct {
	Message string
}

type response struct {
	Emoticons []string      `json:"emoticons"`
	Links     []parser.Link `json:"links"`
	Mentions  []string      `json:"mentions"`
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/entities", Entities).Methods("POST")
	// TODO(jsedlacek): make port configurable via command-line flag (or however they prefer to do their configuration)
	log.Fatal(http.ListenAndServe(":8000", router))
}

func Entities(w http.ResponseWriter, r *http.Request) {
	var req request
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	json.Unmarshal(b, &req)
	msg := req.Message
	emoticons := parser.Emoticons(msg)
	mentions := parser.Mentions(msg)
	links, err := parser.Links(msg)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	resp := response{Emoticons: emoticons, Mentions: mentions, Links: links}
	// TODO(jarek): cache encoder instead of re-creating it?
	json.NewEncoder(w).Encode(resp)
}
