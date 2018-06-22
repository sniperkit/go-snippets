package sessionstore

import (
	"github.com/asim/go-micro/store"
	"github.com/gorilla/sessions"
)

func New() (sessions.Store, error) {
	secretKey, err := store.Get("session/cookie/secret")
	if err != nil {
		return nil, err
	}

	return sessions.NewCookieStore(secretKey.Value()), nil
}
