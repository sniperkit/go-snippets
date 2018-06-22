// Package httputil holds common usefull functions for http
package httputil

import (
	"fmt"
	"net/http"
)

// WriteResponse is a common way to write plain-text responses
func WriteResponse(response interface{}, code int, r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(code)
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	fmt.Fprintf(w, "%d %q", code, response)
}
