package main

import (
	"fmt"
	"net/http"
)

func serveError(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "error.html", struct{}{})
	if err != nil {
		fmt.Println("Error executing template:", err)
	}
}
