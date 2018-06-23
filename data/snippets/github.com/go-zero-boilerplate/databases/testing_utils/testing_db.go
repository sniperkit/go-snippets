package testing_utils

import (
	"os"
	"strings"

	"github.com/jmoiron/sqlx"
)

const (
	TESTS_SAFE_TABLE_NAME = "tests_safety_table"
)

func GetTestingDb(dummyTableNameToCreate string) (*sqlx.DB, error) {
	env_driver := "GODB_TESTS_DRIVER"
	dbDriver := os.Getenv(env_driver)

	env_datasource := "GODB_TESTS_DATASOURCE"
	dataSource := os.Getenv(env_datasource)

	if strings.TrimSpace(dbDriver) == "" {
		panic("Please set the " + env_driver + " environment variable to the driver to use, like mysql, sqlite, etc")
	}
	if strings.TrimSpace(dataSource) == "" {
		panic("Please set the " + env_datasource + " environment variable to the datasource of the driver")
	}

	db, err := sqlx.Connect(dbDriver, dataSource)
	if err != nil {
		return nil, err
	}

	type tmpStruct struct{}
	tmpList := []*tmpStruct{}
	err = db.Select(&tmpList, `SELECT * FROM `+TESTS_SAFE_TABLE_NAME)
	if err != nil {
		return nil, err
	}

	if strings.TrimSpace(dummyTableNameToCreate) != "" {
		_, err = db.Exec(`DROP TABLE IF EXISTS ` + dummyTableNameToCreate)
		if err != nil {
			return nil, err
		}

		_, err = db.Exec(`        
        CREATE TABLE ` + dummyTableNameToCreate + ` (
             id    INTEGER PRIMARY KEY AUTO_INCREMENT
            ,name  VARCHAR(200)
            ,age   INTEGER
        )
        `)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
