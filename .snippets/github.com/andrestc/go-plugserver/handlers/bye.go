package main

import "net/http"

func Handler() (string, func(w http.ResponseWriter, r *http.Request)) {
	return "/bye", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Bye World!"))
	}
}
