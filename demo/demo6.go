/*
* @Author: wuhailin
* @Date:   2017-10-10 08:59:23
* @Last Modified by:   wuhailin
* @Last Modified time: 2017-10-10 09:18:13
 */

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var db *sql.DB

func init() {
	var err error
	driver := `mysql`
	//ifs-test config
	dataSourceName := `ifs_ro:dYDIiSFY@tcp(192.168.7.165)/ifs?charset=utf8`

	//localhost config
	//dataSourceName = `root:123456@tcp(127.0.0.1)/ifs?charset=utf8`

	db, err = sql.Open(driver, dataSourceName)
	if nil != err {
		log.Fatalln(err)
	}
	db.SetMaxOpenConns(10)
}

func main() {
	defer db.Close()
	showTables("mock_model")
}

func showTables(like string) (tables []Table) {
	var table string
	query := `SHOW TABLES`
	if like != "" {
		query = "SHOW TABLES LIKE '" + like + "'"
	}
	rows, err := db.Query(query)
	if nil != err {
		log.Fatalln(err)
	}

	for rows.Next() {
		if err := rows.Scan(&table); nil != err {
			log.Println(err)
		}
		t := Table{name: table}
		tables = append(tables, t)
	}
	return
}

type Table struct {
	name, dll string
	columns   []map[string]string
	total     int
}

func (t *Table) ShowDll() string {
	if "" != t.dll {
		return t.dll
	}
	var name, dll string
	query := "SHOW CREATE TABLE " + t.name
	if err := db.QueryRow(query).Scan(&name, &dll); nil != err {
		log.Fatalln(err)
	}
	t.dll = dll
	return dll
}

func (t *Table) Desc() []map[string]string {
	if nil != t.columns {
		return t.columns
	}
	query := "DESC " + t.name
	rows, err := db.Query(query)
	defer rows.Close()
	if nil != err {
		log.Fatalln(err)
	}
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
		t.columns = append(t.columns, column)
	}
	return t.columns
}

func (t *Table) Total() int {
	if 0 < t.total {
		return t.total
	}
	var total int
	query := "SELECT COUNT(1) FROM " + t.name
	if err := db.QueryRow(query).Scan(&total); nil != err {
		log.Fatalln(err)
	}
	t.total = total
	return total
}
