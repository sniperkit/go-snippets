package sqlimporter

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/jmoiron/sqlx"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

func createDSN() string {
	return ""
}

func connect(driver, dsn string) (*sqlx.DB, error) {
	db, err := sqlx.Open(driver, dsn)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

const DBNameDefault = "SQL_IMPORTER_DB_"

// CreateDatabaseAndImport used to create database
// and import all queries located in a directories
func CreateDB(driver, dsn string) (*sqlx.DB, error, func() error) {
	defaultDrop := func() error {
		return nil
	}
	db, err := connect(driver, dsn)
	if err != nil {
		return nil, err, defaultDrop
	}

	// create a new database
	// database name is always a random name
	unix := time.Now().Unix()
	randSource := rand.NewSource(unix)
	r := rand.New(randSource)
	dbName := DBNameDefault + strconv.Itoa(r.Int())
	// TODO: separate this, this is a dialect and might be not the same with other db
	createDBQuery := fmt.Sprintf("%s %s", getDialect(driver, "create"), dbName)
	// exec create new b
	_, err = db.Exec(createDBQuery)
	if err != nil {
		return nil, err, defaultDrop
	}

	// use new db
	useDatabaseQuery := fmt.Sprintf("USE %s", dbName)
	_, err = db.Exec(useDatabaseQuery)
	if err != nil {
		return nil, err, defaultDrop
	}
	return db, nil, func() error {
		_, err := db.Exec(fmt.Sprintf("%s %s", getDialect(driver, "drop"), dbName))
		if err != nil {
			return err
		}
		return db.Close()
	}
}

// ImportSchemaFromFiles
func ImportSchemaFromFiles(db *sqlx.DB, filepath string) error {
	files, err := getFileList(filepath)
	if err != nil {
		return err
	}

	// an sql file will be executed as one batch of transaction
	for _, file := range files {
		sqlContents, err := parseFiles(file)
		if err != nil {
			return err
		}
		// end if empty
		if len(sqlContents) == 0 {
			return nil
		}

		tx, err := db.BeginTx(context.TODO(), nil)
		if err != nil {
			return err
		}

		for key := range sqlContents {
			_, err = tx.ExecContext(context.TODO(), sqlContents[key])
			if err != nil {
				break
			}
		}

		if err != nil {
			errRollback := tx.Rollback()
			if errRollback != nil {
				return fmt.Errorf("Failed to rollback from file %s with error %s", file, errRollback.Error())
			}
			return fmt.Errorf("Failed to execute file %s with error %s", file, err.Error())
		} else {
			err = tx.Commit()
			if err != nil {
				return fmt.Errorf("Failed to commit from file %s with error %s", file, err.Error())
			}
		}
	}
	return nil
}
