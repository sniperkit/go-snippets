package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome, %!", r.URL.Path[1:])
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/about/", about)
	http.ListenAndServe(":8080", nil)
}

type Message struct {
	Text string
}

func about(w http.ResponseWriter, r *http.Request) {

	m := Message{"Welcome to the SandovalEffect API, build v0.0.001.992, 6/22/2015 0340 UTC."}
	b, err := json.Marshal(m)

	if err != nil {
		panic(err)
	}

	w.Write(b)
}
