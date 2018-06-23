package database

import (
	"../../common"
	"database/sql"
	"fmt"
	"log"
)

func FetchAll(rows *sql.Rows) (result []map[string]interface{}) {
	defer rows.Close()
	cols, _ := rows.Columns()
	data := make([]interface{}, len(cols))
	buf := make([]interface{}, len(cols))

	for i := range buf {
		buf[i] = &data[i]
	}
	for rows.Next() {
		if err := rows.Scan(buf...); nil != err {
			log.Fatalln(err)
		}
		column := make(map[string]interface{})
		for k, v := range data {
			column[cols[k]] = v
		}
		result = append(result, column)
	}
	return
}

func Fetch(rows *sql.Rows) (result map[string]string) {
	defer rows.Close()
	cols, _ := rows.Columns()
	data := make([]interface{}, len(cols))
	buf := make([]interface{}, len(cols))

	for i := range buf {
		buf[i] = &data[i]
	}
	for rows.Next() {
		if err := rows.Scan(buf...); nil != err {
			log.Fatalln(err)
		}
		column := make(map[string]string)
		for k, v := range data {
			if nil == v {
				v = ""
			}
			column[cols[k]] = fmt.Sprintf("%s", v)
		}
		result = column
		break
	}
	return
}

func FetchCol(query string, params ...interface{}) (result string) {
	rows, err := common.Query(query, params)
	if err != nil {
		common.Log(err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&result); nil != err {
			log.Fatalln(err)
		}
		break
	}
	return
}
