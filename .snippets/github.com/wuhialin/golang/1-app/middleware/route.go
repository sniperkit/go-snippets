package middleware

import (
	"net/http"
)

func Route(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}
