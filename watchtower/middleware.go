package main

import (
	"net/http"
)

// DefaultHeadersHandler acts as middleware, adding
// headers like content-type and allow-origin to requests.
type DefaultHeadersHandler struct {
	nextHandler http.Handler
}

// NewDefaultHeadersHandler returns a new DefaultHeadersHandler
func NewDefaultHeadersHandler(handler http.Handler) http.Handler {
	return &DefaultHeadersHandler{nextHandler: handler}
}

func (h *DefaultHeadersHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Add default headers
	if origin := r.Header.Get("Origin"); origin != "" {
		w.Header().Set("Access-Control-Allow-Origin", origin)
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token")
	w.Header().Set("X-Server-Name", "Watchtower")

	// Pass the request to the next handler
	h.nextHandler.ServeHTTP(w, r)
}
