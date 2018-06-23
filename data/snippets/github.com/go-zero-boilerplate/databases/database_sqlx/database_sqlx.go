package database_sqlx

import (
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"

	"github.com/go-zero-boilerplate/databases"
)

func NewSqlxDatabase(db *sqlx.DB) databases.Database {
	return &sqlxDatabase{
		db: db,
	}
}

type sqlxDatabase struct {
	db *sqlx.DB
}

func (s *sqlxDatabase) IsErrNoRows(err error) bool {
	return sql.ErrNoRows == err
}

func (s *sqlxDatabase) Select(dest interface{}, query string, args ...interface{}) error {
	return s.db.Select(dest, query, args...)
}

func (s *sqlxDatabase) Exec(query string, args ...interface{}) (databases.ExecResult, error) {
	result, err := s.db.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *sqlxDatabase) QueryRow(query string, args ...interface{}) databases.ResultRow {
	return s.db.QueryRowx(query, args...)
}

func (s *sqlxDatabase) Query(query string, args ...interface{}) (databases.ResultRows, error) {
	return s.db.Queryx(query, args...)
}

func (s *sqlxDatabase) BeginTx() (databases.Database, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	return &sqlxTransaction{
		tx:      tx,
		isOwner: true,
	}, nil
}

func (s *sqlxDatabase) CommitTx() error {
	return errors.New("[databases] Unexpected call to CommitTx (non-transaction) Database instance")
}

func (s *sqlxDatabase) RollbackTx() error {
	return errors.New("[databases] Unexpected call to RollbackTx (non-transaction) Database instance")
}

func (s *sqlxDatabase) DeferredRollbackIfNotHandled() error {
	return errors.New("[databases] Unexpected call to DeferredRollbackIfNotHandled (non-transaction) Database instance")
}
