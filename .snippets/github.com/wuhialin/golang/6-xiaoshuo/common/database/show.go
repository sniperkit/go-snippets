package database

import (
	"database/sql"
	"log"
)

var Show show

func init() {
	Show = show{}
}

type show struct {
	db *sql.DB
}

func (t *show) Tables() (tables []string) {
	rows, err := t.db.Query(`SHOW TABLES`)
	if err != nil {
		log.Println(err)
		return
	}
	for _, data := range FetchAll(rows) {
		for _, table := range data {
			tables = append(tables, table.(string))
		}
	}
	return
}
