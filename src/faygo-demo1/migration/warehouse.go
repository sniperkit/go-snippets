package migration

import (
	"github.com/henrylee2cn/faygo"
	"faygo-demo1/model"
)

func warehouse() {
	var warehouse model.Warehouse
	warehouse.Last(&warehouse)
	rows, err := db.Debug().Table("c_warehouse").Select("id, warehouse_name, warehouse_code").Where("id > ?", warehouse.ID).Rows()
	if err != nil {
		faygo.Error(err)
		return
	}
	for rows.Next() {
		m := model.Warehouse{}
		rows.Scan(&m.ID, &m.Name, &m.Code)
		if err = m.Save(&m).Error; err != nil {
			faygo.Error(err)
			break
		}
	}
}