/*
* @Author: wuhailin
* @Date:   2018-01-20 14:28:44
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-22 15:23:19
 */
package ifs

import (
	"github.com/jinzhu/gorm"
)

func FetchCategory(db *gorm.DB, id uint, callback func(map[string]interface{})) {
	rows, err := db.Table("p_product_catalog").Where("id > ?", id).Rows()
	if err != nil {
		panic(err)
	}
	columns, _ := rows.Columns()
	columnLen := len(columns)
	buf := make([]interface{}, columnLen)
	data := make([]interface{}, columnLen)
	for k := range buf {
		buf[k] = &data[k]
	}
	row := make(map[string]interface{})
	for rows.Next() {
		rows.Scan(buf...)
		for k, v := range data {
			row[columns[k]] = v
		}
		callback(row)
	}
}
