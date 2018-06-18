package main

import (
	"database/sql"
	stderr "errors"
	"log"

	"github.com/albert-widi/go_common/errors"
	"github.com/albert-widi/go_common/logger/go-kit"
)

func main() {
	l := logger.New()
	err := errors.New("Ini error", []string{"ini field1", "ini field2"})
	err2 := stderr.New("Ini error")
	err3 := errors.New(sql.ErrNoRows, errors.Fields{"order_id": 10000, "process": "add_order_id"})
	if errors.Match(err, err2) {
		l.Print("Errornya sama bos")
		log.Print(err.GetMessages())
	} else {
		l.Print("Beda bos")
	}
	l.Errors(err3)
}
