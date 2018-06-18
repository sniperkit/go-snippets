package main

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type movie struct {
	ID          bson.ObjectId `bson:"_id" json:"id"`
	Name        string        `bson:"name" json:"name"`
	CoverImage  string        `bson:"cover_image" json:"coverImage"`
	Description string        `bson:"description" json:"description"`
}

const (
	COLLECTION = "movies"
)

func getMovies(db *mgo.Database) ([]movie, error) {
	var m []movie
	err := db.C(COLLECTION).Find(bson.M{}).All(&m)

	return m, err
}

func getMovie(db *mgo.Database, id string) (movie, error) {
	var m movie
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&m)

	return m, err
}

func (m *movie) createMovie(db *mgo.Database) error {
	err := db.C(COLLECTION).Insert(&m)

	return err
}

func (m *movie) updateMovie(db *mgo.Database) error {
	err := db.C(COLLECTION).UpdateId(m.ID, &m)

	return err
}

func (m *movie) deleteMovie(db *mgo.Database) error {
	err := db.C(COLLECTION).Remove(&m)

	return err
}
