package api

import (
	"fmt"
	"net/http"
)

func (api *Router) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "{\"status\": \"OK\"}")
}
