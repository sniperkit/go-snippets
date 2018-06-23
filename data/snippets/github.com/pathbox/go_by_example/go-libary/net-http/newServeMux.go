package main

import (
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "just for test!")
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("/home")))
	mux.HandleFunc("/test", Test)
	err := http.ListenAndServe(":9999", mux)
	if err != nil {
		fmt.Println(err)
	}
}
