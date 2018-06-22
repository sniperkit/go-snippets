package main

import (
	"fmt"
	"net/http"

	sessions "github.com/goincremental/negroni-sessions"
)

func RequireAuth(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	s := sessions.GetSession(r)

	//not authorized redirect to login
	if s.Get("login") == nil {
		RedirectToLogin(w, r)
	}

	next(w, r)
}

func RedirectToLogin(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/login", http.StatusFound)
}

func RedirectToError(w http.ResponseWriter, r *http.Request, err error) {
	http.Redirect(w, r, fmt.Sprint("/error?err=", err.Error()), http.StatusFound)
}

func ValidateRoute(s string, w http.ResponseWriter, r *http.Request, id string) {
	if id != "" {
		if r.URL.Path != "/"+s+"/"+id {
			NotFound(w, r)
			return
		}
	} else {
		if r.URL.Path != "/"+s {
			NotFound(w, r)
			return
		}
	}
}
