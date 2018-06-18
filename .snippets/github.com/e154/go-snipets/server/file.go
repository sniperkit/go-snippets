package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "applicaiton/zip")
			w.Header().Set("Content-Disposition", "attachment; filename='buf.zip'")
			http.ServeFile(w, r, "/tmp/buf.zip")
		})

	http.ListenAndServe(":8088", nil)
}
