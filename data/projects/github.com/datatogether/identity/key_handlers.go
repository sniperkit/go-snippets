package main

import (
	"github.com/datatogether/identity/user"

	// "encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

func KeysHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "OPTIONS":
		CORSHandler(w, r)
	case "GET":
		UserKeysHandler(w, r)
	case "POST":
		CreateUserKeyHandler(w, r)
	case "DELETE":
		DeleteUserKeyHandler(w, r)
	default:
		ErrRes(w, ErrNotFound)
	}
}

func CreateUserKeyHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	if u == nil {
		ErrRes(w, ErrNotFound)
		return
	}

	req := struct {
		Name string
		Key  string
	}{}

	if isJsonRequest(r) {
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			log.Info(err)
			ErrRes(w, err)
			return
		}
	} else {
		req.Name = r.FormValue("name")
		req.Key = r.FormValue("key")
	}

	key, err := user.CreateKey(appDB, u, req.Name, []byte(req.Key))
	if err != nil {
		log.Info(err)
		ErrRes(w, err)
		return
	}

	Res(w, true, key)
}

func UserKeysHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	if u == nil {
		ErrRes(w, ErrNotFound)
		return
	}

	keys, err := u.Keys(appDB)
	if err != nil {
		log.Info(err)
		ErrRes(w, err)
		return
	}

	Res(w, true, keys)
}

func DeleteUserKeyHandler(w http.ResponseWriter, r *http.Request) {
	u := sessionUser(r)
	if u == nil {
		ErrRes(w, ErrNotFound)
		return
	}

	// TODO - finish
	ErrRes(w, fmt.Errorf("deleting keys is not yet implemented"))
	return

	// shaStr := stringParam(ctx, "id")
	// sha256 := [32]byte{}
	// sha, err := hex.DecodeString(shaStr)
	// if err != nil {
	// 	ErrRes(w, ErrNotFound)
	// 	return
	// }
	// for i, b := range sha {
	// 	sha256[i] = byte(b)
	// }

	// log.Info(shaStr, sha256)
	// key := &UserKey{Sha256: sha256}
	// if err := key.Delete(appDB); err != nil {
	// 	log.Info(err)
	// 	ErrRes(w, err)
	// 	return
	// }

	// Res(w, key)
}
