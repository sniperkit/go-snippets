package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/datatogether/sqlutil"
	"github.com/gchaincl/dotsql"
	_ "github.com/lib/pq"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()

	var err error
	cfg, err = initConfig(TEST_MODE) // make sure we read env in test mode
	if err != nil {
		panic(err)
	}

	teardown := setupTestDatabase()

	retCode := m.Run()
	teardown()
	os.Exit(retCode)
}

func setupTestDatabase() func() {
	var err error
	appDB, err = sqlutil.SetupConnection("postgres", cfg.PostgresDbUrl)
	if err != nil {
		appDB.Close()
		log.Panicln(err)
	}

	teardown, err := initializeAppSchema(appDB)
	if err != nil {
		log.Panicln(err.Error())
	}

	if err := resetTestData(appDB, "users", "reset_tokens", "keys"); err != nil {
		log.Panicln(err.Error())
	}

	return teardown
}

// WARNING - THIS ZAPS WHATEVER DB IT'S GIVEN. DO NOT CALL THIS SHIT.
// used for testing only, returns a teardown func
func initializeAppSchema(db *sql.DB) (func(), error) {
	// if cfg.Mode != TEST_MODE {
	// 	return errors.New("attempted to initialize schema while not in test mode")
	// }

	schema, err := dotsql.LoadFromFile("sql/schema.sql")
	if err != nil {
		return nil, err
	}

	for _, cmd := range []string{
		"drop-all",
		"create-users",
		"create-reset_tokens",
		"create-keys",
	} {
		if _, err := schema.Exec(db, cmd); err != nil {
			log.Println(cmd, "error:", err)
			return nil, err
		}
	}

	teardown := func() {
		if _, err := schema.Exec(db, "drop-all"); err != nil {
			log.Panicln(err.Error())
		}
	}

	return teardown, nil
}

// drops test data tables & re-inserts base data from sql/test_data.sql, based on
// passed in table names
func resetTestData(db *sql.DB, tables ...string) error {
	schema, err := dotsql.LoadFromFile("sql/test_data.sql")
	if err != nil {
		return err
	}
	for _, t := range tables {
		if _, err := schema.Exec(db, fmt.Sprintf("delete-%s", t)); err != nil {
			return fmt.Errorf("error deleting %s: %s", t, err.Error())
		}
		if _, err := schema.Exec(db, fmt.Sprintf("insert-%s", t)); err != nil {
			return fmt.Errorf("error inserting %s: %s", t, err.Error())
		}
	}
	return nil
}
