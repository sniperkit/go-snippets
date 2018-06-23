package data

import (
	"fmt"

	mgo "gopkg.in/mgo.v2"
)

// DB to encapsulate the mongodb access and provide helpers for collections.
// TODO: Handle session copying, connection pooling, etc.
type DB struct {
	session  *mgo.Session
	Sessions *mgo.Collection
}

// NewDB creates a new database which connects to mongodb.
func NewDB() *DB {
	ip := "127.0.0.1"
	s, err := mgo.Dial(ip)

	if err != nil {
		fmt.Printf("Error: Unable to connect to MongoDB (%s)\n", ip)
	}

	db := &DB{
		session:  s,
		Sessions: s.DB("gg").C("sessions"),
	}

	return db
}
