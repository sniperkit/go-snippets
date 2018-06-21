/*
* @Author: wuhialin
* @Date:   2018-01-17 18:26:32
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-19 16:42:49
 */
package ifs

import "github.com/jinzhu/gorm"

func FetchWarehouse(db *gorm.DB, id uint, callback func(map[string]interface{})) {
	rows, err := db.Table("c_warehouse").Select(`id, warehouse_name, warehouse_code, auto_send, type, warehouse_type, is_oversea, dept_id, special_data_examine, special_demand_examine, demand_examine, stock_yun_code, demand_auto_pass`).Where("id > ?", id).Rows()
	if err != nil {
		panic(err)
	}
	var name, sn, stockYunCode string
	var autoSend, Type, warehouseType, isOversea, specialDataExamine, specialDemandExamine, demandExamine, demandAutoPass uint8
	var deptId uint
	row := make(map[string]interface{})
	for rows.Next() {
		rows.Scan(&id, &name, &sn, &autoSend, &Type, &warehouseType, &isOversea, &deptId, &specialDataExamine, &specialDemandExamine, &demandExamine, &stockYunCode, &demandAutoPass)
		row["id"] = id
		row["name"] = name
		row["sn"] = sn
		row["autoSend"] = autoSend
		row["Type"] = Type
		row["warehouseType"] = warehouseType
		row["isOversea"] = isOversea
		row["deptId"] = deptId
		row["specialDataExamine"] = specialDataExamine
		row["specialDemandExamine"] = specialDemandExamine
		row["demandExamine"] = demandExamine
		row["stockYunCode"] = stockYunCode
		row["demandAutoPass"] = demandAutoPass
		callback(row)
	}
}

func FetchWarehouseDept(db *gorm.DB, id uint, callback func(map[string]interface{})) {
	rows, err := db.Table("c_warehouse_dept").Select(`id, department_name, department_sn, status`).Where("id > ?", id).Rows()
	if err != nil {
		panic(err)
	}
	var status uint8
	var name, sn string
	row := make(map[string]interface{})
	for rows.Next() {
		rows.Scan(&id, &name, &sn, &status)
		row["id"] = id
		row["name"] = name
		row["sn"] = sn
		row["status"] = status
		callback(row)
	}
}
