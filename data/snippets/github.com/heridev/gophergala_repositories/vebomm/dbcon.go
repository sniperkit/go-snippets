package main

import (
	"database/sql"
	"fmt"
	"log"

	. "github.com/gophergala/vebomm/core"

	"github.com/go-gorp/gorp"
	_ "github.com/lib/pq"
)

const (
	DbType = "postgres"
)

var (
	gDb *gorp.DbMap
)

func Db() *gorp.DbMap {
	if gDb == nil {
		panic("Database service uninitialized.")
	}
	return gDb
}

type DbConfig struct {
	Database string
	User     string
	Password string
	Host     string
	Port     int
}

func loadDb(conf DbConfig) *sql.DB {
	db, err := sql.Open(DbType,
		fmt.Sprintf("%v://%v:%v@%v:%v/%v",
			DbType,
			conf.User,
			conf.Password,
			conf.Host,
			conf.Port,
			conf.Database,
		))
	if err != nil {
		log.Println("Cannot open database: " + err.Error())
	}

	return db
}

func InitDb(conf DbConfig) {
	db := loadDb(conf)
	if db != nil {
		constructDb(db)
	}
}

func constructDb(db *sql.DB) {
	// construct a gorp DbMap
	dbMap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}

	dbMap.AddTableWithName(User{}, "users").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	err := dbMap.CreateTablesIfNotExists()
	if err != nil {
		log.Println("Failed creating tables. " + err.Error())
	}

	gDb = dbMap
}
