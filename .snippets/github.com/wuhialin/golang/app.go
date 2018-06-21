package main

import (
	"io"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		log.Println(req.RemoteAddr)
		io.WriteString(w, "hello world!\n")
	})
	http.ListenAndServe(":8080", nil)
}
