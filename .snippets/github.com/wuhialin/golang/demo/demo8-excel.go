package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tealeg/xlsx" //数据一大，内存使用太严重
	"log"
)

func main() {
	driver := "root:123456@tcp(127.0.0.1)/ifs?charset=utf8"
	db, err := sql.Open("mysql", driver)
	if nil != err {
		log.Fatalln(err)
	}
	defer db.Close()
	rows, _ := db.Query("SELECT id, goods_sn, warehouse_code, a_story_quantity, b_story_quantity, c_story_quantity, create_date, add_time FROM b_story LIMIT 1000")
	defer rows.Close()
	cols, _ := rows.Columns()
	data := make([]interface{}, len(cols))
	buf := make([]interface{}, len(cols))
	for i := range buf {
		buf[i] = &data[i]
	}

	var file *xlsx.File
	var sheet *xlsx.Sheet
	var row *xlsx.Row
	var cell *xlsx.Cell
	file = xlsx.NewFile()
	sheet, err = file.AddSheet("Sheet1")
	if nil != err {
		log.Fatalln(err)
	}

	for rows.Next() {
		rows.Scan(buf...)
		row = sheet.AddRow()
		for _, v := range data {
			cell = row.AddCell()
			cell.Value = fmt.Sprintf(`%s`, v)
		}
	}
	err = file.Save("story.xlsx")
	if nil != err {
		log.Fatalln(err)
	}
}
