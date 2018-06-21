package common

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"strings"
)

func db() *sql.DB {
	name := "postgres"
	//c := config()["database"]
	source := "postgres://postgres:123456@localhost/postgres?sslmode=disable"
	db, err := sql.Open(name, source)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(1000)
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	return db
}

func Query(query string, params ...interface{}) (rows *sql.Rows, err error) {
	rows, err = connDb.Query(query, params...)
	return
}

func Exec(query string, params ...interface{}) (result sql.Result, err error) {
	result, err = connDb.Exec(query, params...)
	return
}

func buildInsertSql(table string, fields []string, params ...interface{}) (sqls []string) {
	fieldLen := len(fields)
	if fieldLen == 0 || len(params) == 0 {
		return
	}
	sqls = append(sqls, "INSERT INTO", table)
	tmpField := make([]string, fieldLen)
	for k, v := range fields {
		tmpField[k] = v
	}
	sqls = append(sqls, `(`+strings.Join(tmpField, ", ")+`)`, "VALUES")
	key := 0
	var values []string
	for k := range params {
		key = k % fieldLen
		tmpField[key] = fmt.Sprintf("$%d", k+1)
		if key = key + 1; key == fieldLen {
			//每一行分组
			values = append(values, `(`+strings.Join(tmpField, ", ")+`)`)
		}
	}
	sqls = append(sqls, strings.Join(values, ", "))
	return
}

func Insert(table string, fields []string, params ...interface{}) (id int64, err error) {
	sqls := buildInsertSql(table, fields, params...)
	sqls = append(sqls, "RETURNING id")
	db().QueryRow(strings.Join(sqls, " "), params...).Scan(&id)
	return
}

func BatchInsert(table string, fields []string, params ...interface{}) (result sql.Result, err error) {
	sqls := buildInsertSql(table, fields, params...)
	sqls = append(sqls, "ON CONFLICT DO NOTHING")
	return Exec(strings.Join(sqls, " "), params...)
}
