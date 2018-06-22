package main

import (
	"os"

	"gopkg.in/mgo.v2"
)

// Generic insertion statement
func Insert(collectionName string, object interface{}) error {
	mongoUri := connectionString()
	sess, err := mgo.Dial(mongoUri)
	if err != nil {
		return err
	}
	defer sess.Close()

	sess.SetSafe(&mgo.Safe{})
	c := sess.DB("timeywimey").C(collectionName)
	err = c.Insert(object)
	if err != nil {
		return err
	}

	return nil
}

// Get the MongoDB connection string
func connectionString() string {
	env := os.Getenv("TW_MONGO_URL")
	if env == "" {
		return "localhost:27017"
	}
	return env
}
