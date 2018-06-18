package main

import (
	"github.com/albert-widi/go_common/logger/go-kit"
	"github.com/albert-widi/go_common/sqlimporter"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	l := logger.New()
	db, err, drop := sqlimporter.CreateDB("mysql", "root:@tcp(127.0.0.1:3306)/?parseTime=true")
	if err != nil {
		l.Print(err)
	}
	l.Print("Database created")
	if err := db.Ping(); err != nil {
		l.Print("Database is not reachable")
	}

	err = sqlimporter.ImportSchemaFromFiles(db, "../files")
	if err != nil {
		l.Print(err)
	}
	defer drop()
}
