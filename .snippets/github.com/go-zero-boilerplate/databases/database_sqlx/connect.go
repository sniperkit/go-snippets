package database_sqlx

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

func Connect(driverName, dataSourceName string) (*sqlx.DB, error) {
	return sqlx.Connect(driverName, dataSourceName)
}

func Open(driverName, dataSourceName string) (*sqlx.DB, error) {
	return sqlx.Open(driverName, dataSourceName)
}

func NewDb(db *sql.DB, driverName string) *sqlx.DB {
	return sqlx.NewDb(db, driverName)
}
