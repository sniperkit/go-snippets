package rethinkdb

import (
	"github.com/asim/go-micro/store"
	r "github.com/dancannon/gorethink"
)

var session *r.Session

func Session() (*r.Session, error) {
	var err error

	if session == nil {
		var address, database string

		item, err := store.Get("rethinkdb/address")
		if err != nil {
			address = "localhost:28015"
		} else {
			address = string(item.Value())
		}
		item, err = store.Get("rethinkdb/database")
		if err != nil {
			database = "gevvent"
		} else {
			database = string(item.Value())
		}

		session, err = r.Connect(r.ConnectOpts{
			Address:  address,
			Database: database,
		})
		if err != nil {
			session = nil
		}
	}

	return session, err
}
