package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	listen := "127.0.0.1:3000"

	fmt.Println("Starting Webserver on", listen)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Go REST Service (Query: %q)", html.EscapeString(r.URL.Path))
	})
	log.Fatal(http.ListenAndServe(listen, nil))
}
