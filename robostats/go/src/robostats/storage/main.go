// Package storage provides functions to access permanent storage resources.
package storage

import (
	"log"
	"upper.io/db"
	"upper.io/db/mongo"
)

// Settings returns settings for connecting to the database.
func Settings() db.ConnectionURL {
	// TODO: Move this stuff to a file.
	settings := mongo.ConnectionURL{
		Address:  db.Host("mongo-server"),
		Database: "robostats",
	}
	return settings
}

// C returns a collection or exits on error.
func C(name string) db.Collection {
	col, err := DB().Collection(name)
	if err != nil {
		if err != db.ErrCollectionDoesNotExist {
			log.Fatalf("conn.C(%s): %q", name, err)
		}
	}
	return col
}

// DB returns a database session or exists on error.
func DB() db.Database {
	sess, err := db.Open(mongo.Adapter, Settings())
	if err != nil {
		log.Fatalf("conn.DB(): %q", err)
	}
	return sess
}
