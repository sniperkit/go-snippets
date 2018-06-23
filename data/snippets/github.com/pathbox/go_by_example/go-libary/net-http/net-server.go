package main

import (
	"fmt"
	"net/http"
)

func Test(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "just for test!")
}
func main() {
	newserver := http.Server{
		Addr:         ":9999",
		ReadTimeout:  0,
		WriteTimeout: 0,
	}
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./")))
	mux.HandleFunc("/test", Test)
	newserver.Handler = mux
	err := newserver.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(err)
	// err := http.ListenAndServe(":9999", mux)
	// if err != nil {
	//     fmt.Println(err)
	// }
}
