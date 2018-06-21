package common

import "database/sql"

var openDb *sql.DB

func DB() *sql.DB {
	if openDb != nil {
		return openDb
	}
	name := "mysql"
	source := ""
	db, err := sql.Open(name, source)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(100)
	openDb = db
	return db
}
