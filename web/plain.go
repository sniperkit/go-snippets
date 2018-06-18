package main

import (
	"net/http"
)

func plain(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func main() {
	http.HandleFunc("/", plain)
	http.ListenAndServe(":3000", nil)
}

