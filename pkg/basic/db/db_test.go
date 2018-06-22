package db

import (
	"context"
	"database/sql"
	"log"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDB(t *testing.T) {
	db, err := sql.Open("mysql", "http://localhost:9001")
	log.Printf("db: %v, err: %v", db, err)
	log.Println(db.Conn(context.Background()))
}
