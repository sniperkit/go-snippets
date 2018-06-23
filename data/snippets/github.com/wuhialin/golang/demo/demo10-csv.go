package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

func main() {
	driver := "root:123456@tcp(127.0.0.1)/ifs?charset=utf8"
	db, err := sql.Open("mysql", driver)
	if nil != err {
		log.Fatalln(err)
	}
	defer db.Close()
	rows, _ := db.Query("SELECT * FROM b_story")
	defer rows.Close()
	cols, _ := rows.Columns()
	data := make([]interface{}, len(cols))
	buf := make([]interface{}, len(cols))
	for i := range buf {
		buf[i] = &data[i]
	}
	f, _ := os.OpenFile("story.csv", os.O_WRONLY|os.O_APPEND|os.O_CREATE, os.ModePerm)
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	for rows.Next() {
		rows.Scan(buf...)
		var row []string
		for _, v := range data {
			row = append(row, fmt.Sprintf(`%s`, v))
		}
		w.Write(row)
	}
}
