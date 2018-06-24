// abstractions for working with postgres databases
package main

import (
	"database/sql"
	"fmt"
	"github.com/gchaincl/dotsql"
	"time"

	_ "github.com/lib/pq"
)

// this interface unifies both *sql.Row & *sql.Rows
type sqlScannable interface {
	Scan(...interface{}) error
}

// sqlQuerable unifies both *sql.DB & *sql.Tx for querying purposes
type sqlQueryable interface {
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type sqlExecable interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

type sqlQueryExecable interface {
	sqlQueryable
	sqlExecable
}

func connectToAppDb() {
	var err error
	fmt.Println("connecting to db")
	for i := 0; i < 1000; i++ {
		appDB, err = SetupConnection(cfg.PostgresDbUrl)
		if err != nil {
			fmt.Println(err)
			time.Sleep(time.Second)
			continue
		}
		fmt.Println("connected to db")
		if err := initializeDatabase(appDB); err != nil {
			fmt.Println(err.Error())
		}
		break
	}
}

// Sets up a connection with a given postgres db connection string
func SetupConnection(connString string) (db *sql.DB, err error) {
	db, err = sql.Open("postgres", connString)
	if err != nil {
		return
	}
	if err = db.Ping(); err != nil {
		return
	}
	return
}

// WARNING - THIS ZAPS WHATEVER DB IT'S GIVEN. DO NOT CALL THIS SHIT.
// used for testing only, returns a teardown func
func initializeDatabase(db *sql.DB) error {
	var err error

	// test query to check for database schema existence
	var exists bool
	if err = db.QueryRow("select exists(select * from primers limit 1)").Scan(&exists); err == nil {
		return nil
	}

	fmt.Println("initializing database with base test data")

	schema, err := dotsql.LoadFromFile(packagePath("/sql/schema.sql"))
	if err != nil {
		return err
	}

	for _, cmd := range []string{
		"drop-all",
		"create-primers",
		"create-sources",
		"create-urls",
		"create-links",
		"create-metadata",
		"create-snapshots",
		"create-collections",
		"create-archive_requests",
		"create-uncrawlables",
		"create-collection_items",
	} {
		if _, err := schema.Exec(db, cmd); err != nil {
			fmt.Println(cmd, "error:", err)
			return err
		}
	}

	if err := insertTestData(
		appDB,
		"primers",
		"sources",
		"urls",
		"links",
		"metadata",
		"snapshots",
		"collections",
		"archive_requests",
		"uncrawlables",
		"collection_items",
	); err != nil {
		return err
	}

	return nil
}

// drops test data tables & re-inserts base data from sql/test_data.sql, based on
// passed in table names
func insertTestData(db *sql.DB, tables ...string) error {
	schema, err := dotsql.LoadFromFile(packagePath("sql/test_data.sql"))
	if err != nil {
		return err
	}
	for _, t := range tables {
		if _, err := schema.Exec(db, fmt.Sprintf("insert-%s", t)); err != nil {
			err = fmt.Errorf("error insert-%s: %s", t, err.Error())
			return err
		}
	}
	return nil
}
