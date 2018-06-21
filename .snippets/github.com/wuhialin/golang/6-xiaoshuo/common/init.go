package common

import (
	"database/sql"
	"log"
)

var connDb *sql.DB

const (
	PRINT_LOG = true
)

func init() {
	connDb = db()
	err := connDb.Ping()
	if err != nil {
		panic(err)
	}
}

func Log(args ...interface{}) {
	if PRINT_LOG {
		log.Println(args...)
	}
}
