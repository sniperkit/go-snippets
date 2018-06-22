package godb

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title"}).
		AddRow(1, "user1").
		AddRow(2, "user2")

	mock.ExpectQuery("^SELECT (.+) FROM user$").WillReturnRows(rows)

	getUser(db)
}
