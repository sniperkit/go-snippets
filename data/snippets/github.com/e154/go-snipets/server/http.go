package main

import (
	"fmt"
	"log"
	"net/http"
)

const listenAddr = "localhost:4000"

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, web")
}

func main() {
	http.HandleFunc("/", handler)
	err := http.ListenAndServe(listenAddr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
