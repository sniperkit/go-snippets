package main

import "net/http"

func Handler() (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}
}
