package main

import (
	"encoding/json"
	"fmt"
	"github.com/datatogether/identity/jwt"
	"github.com/datatogether/identity/user"
	"io"
	"net/http"
)

func JwtHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET", "POST":
		JwtTokenHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

// reads the form values, checks them and creates the token
func JwtTokenHandler(w http.ResponseWriter, r *http.Request) {
	var session *user.User

	if u, ok := r.Context().Value("user").(*user.User); ok {
		session = u
	} else {
		u, err := reqLoginUser(r)
		if err != nil {
			ErrRes(w, err)
			return
		}
		session = u
	}

	tokenString, err := jwt.CreateToken(session)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintln(w, "Sorry, error while Signing Token!")
		log.Info("Token Signing error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/jwt")
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, tokenString)
}

// reqLoginUser facilitates logging in via json or header values
// TODO - should this go in the user package?
func reqLoginUser(r *http.Request) (*user.User, error) {
	var login struct {
		Username string
		Password string
	}

	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
			return nil, NewFmtError(http.StatusBadRequest, "error parsing json: '%s'", err.Error())
		}
	} else {
		// default to using form values
		login.Username = r.Header.Get("username")
		login.Password = r.Header.Get("password")
	}

	u, err := user.AuthenticateUser(appDB, login.Username, login.Password)
	if err != nil {
		return nil, ErrInvalidUserNamePasswordCombo
	}
	return u, nil
}
