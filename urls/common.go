package urls

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/unrolled/render"
)

func Render() *render.Render {
	o := render.Options{}
	o.IndentJSON = true
	o.IndentXML = true
	o.Extensions = []string{"html"}
	return render.New(o)
}

var openDb *sql.DB

func DB() *sql.DB {
	driverName := "mysql"
	driver := "root:123456@tcp(127.0.0.1)/ifs?charset=utf8"
	db, err := sql.Open(driverName, driver)
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(20)
	openDb = db
	return db
}

func FetchAll(rows *sql.Rows) (result []map[string]interface{}) {
	cols, _ := rows.Columns()
	data := make([]interface{}, len(cols))
	buf := make([]interface{}, len(cols))
	for i := range data {
		buf[i] = &data[i]
	}

	for rows.Next() {
		rows.Scan(buf...)
		row := make(map[string]interface{}, len(cols))
		for k, v := range data {
			if nil == v {
				v = ""
			}
			v = fmt.Sprintf("%s", v)
			row[cols[k]] = v
		}
		result = append(result, row)
	}
	return
}
