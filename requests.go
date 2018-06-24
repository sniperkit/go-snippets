package main

import (
	"net/http"
	"strconv"
)

const (
	httpApiInterface = "http-api"
	// sshApiInterface  = "ssh-api"
	// rpcApiInterface  = "rpc-api"
)

type Request interface {
	Exec() (interface{}, error)
}

func ExecRequest(w http.ResponseWriter, envelope bool, req Request) {
	if res, err := req.Exec(); err != nil {
		ErrRes(w, err)
	} else {
		Res(w, envelope, res)
	}
}

func isJsonRequest(r *http.Request) bool {
	// TODO - make this more robust
	return r.Header.Get("Content-Type") == "application/json"
}

// get a boolean query value from an http request
func reqParamBool(key string, r *http.Request) (bool, error) {
	str := r.FormValue(key)
	if str == "true" {
		return true, nil
	}
	return false, nil
}

// get an int query value from an http request
func reqParamInt(key string, r *http.Request) (int, error) {
	str := r.FormValue(key)
	i, err := strconv.ParseInt(str, 10, 0)
	return int(i), err
}
