package main

import (
	"flag"
	"log"
	"net/http"
)

var httpAddr = flag.String("http", ":7000", "HTTP service address")

func init() {
	// Set up the HTTP handler in init (not main) so we can test it. (This main
	// doesn't run when testing.)
	http.HandleFunc("/", home)
}

func main() {
	flag.Parse()
	if err := http.ListenAndServe(*httpAddr, nil); err != nil {
		log.Fatal(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	// FILLIN
}
