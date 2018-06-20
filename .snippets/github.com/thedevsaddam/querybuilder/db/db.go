//author: Saddam H
//author email: thedevsaddam@gmail.com
//package name: querybuilder

package db

import (
	"errors"
	"sync"
)

// DB ...
type DB interface {
	Table(name string)     //set table name
	Query(q string) string //perform raw query on the table
	Find(id int) string    //find a record by it's primary id
	Get() string           //fetch all the records from the table
}

const (
	_ = iota
	// MYSQL ...
	MYSQL
	// SQLITE ...
	SQLITE
	// MONGODB ...
	MONGODB
)

var db DB          //will be used as a singleton db object
var once sync.Once //make thread safe singleton

// New Return a db singleton object
func New(driver int) (DB, error) {
	var err error
	once.Do(func() {
		db, err = databaseFactory(driver)
	})
	return db, err
}

//Database factory
//Return a DB for general purpose
func databaseFactory(driver int) (DB, error) {
	switch driver {
	case MYSQL:
		return new(Mysql), nil
	case SQLITE:
		return new(Sqlite), nil
	case MONGODB:
		return new(MongoDB), nil
	default:
		return nil, errors.New("unsupported storage driver")
	}
}
