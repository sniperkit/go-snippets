package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {

	port := os.Getenv("PORT")

	// home route
	http.HandleFunc("/", homeHandler)

	// static route
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// webapp route
	http.HandleFunc("/app/", appHandler)

	// websocket

	// listen
	fmt.Println("Listening.. on PORT:", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}

}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "home/")
}

func appHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, r.URL.Path[1:])
}
