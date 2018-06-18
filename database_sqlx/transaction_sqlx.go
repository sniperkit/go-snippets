package database_sqlx

import (
	"database/sql"

	"github.com/go-zero-boilerplate/databases"
	"github.com/jmoiron/sqlx"
)

type sqlxTransaction struct {
	tx      *sqlx.Tx
	isOwner bool
	handled bool
}

func (s *sqlxTransaction) IsErrNoRows(err error) bool {
	return sql.ErrNoRows == err
}

func (s *sqlxTransaction) Select(dest interface{}, query string, args ...interface{}) error {
	return s.tx.Select(dest, query, args...)
}

func (s *sqlxTransaction) Exec(query string, args ...interface{}) (databases.ExecResult, error) {
	result, err := s.tx.Exec(query, args...)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (s *sqlxTransaction) QueryRow(query string, args ...interface{}) databases.ResultRow {
	return s.tx.QueryRowx(query, args...)
}

func (s *sqlxTransaction) Query(query string, args ...interface{}) (databases.ResultRows, error) {
	return s.tx.Queryx(query, args...)
}

func (s *sqlxTransaction) BeginTx() (databases.Database, error) {
	//A transaction already started, so just return same instance but with "isOwner" = FALSE
	return &sqlxTransaction{
		tx:      s.tx,
		isOwner: false,
	}, nil
}

func (s *sqlxTransaction) CommitTx() error {
	if !s.isOwner {
		//TODO: Is it fine to just do nothing if we are not the "owner"?
		return nil
	}

	s.handled = true
	return s.tx.Commit()
}

func (s *sqlxTransaction) RollbackTx() error {
	if !s.isOwner {
		//TODO: Is it fine to just do nothing if we are not the "owner"?
		return nil
	}

	s.handled = true
	return s.tx.Rollback()
}

func (s *sqlxTransaction) DeferredRollbackIfNotHandled() error {
	if !s.isOwner {
		//TODO: Is it fine to just do nothing if we are not the "owner"?
		return nil
	}

	if s.handled {
		return nil
	}
	return s.RollbackTx()
}
