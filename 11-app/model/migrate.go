/*
* @Author: wuhailin
* @Date:   2018-01-20 14:55:20
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-22 15:12:20
 */
package model

import (
	"./ifs"
	"fmt"
	"log"
)

func migrateIfs() {
	var err error
	warehouse := Warehouse{}
	err = db.Last(&warehouse).Error
	ifs.FetchWarehouse(mysqlDb, warehouse.ID, func(row map[string]interface{}) {
		warehouse := Warehouse{}
		warehouse.ID = row["id"].(uint)
		warehouse.Name = row["name"].(string)
		warehouse.Code = row["sn"].(string)
		warehouse.AutoSend = row["autoSend"].(uint8)
		warehouse.DemandAutoPass = row["demandAutoPass"].(uint8)
		warehouse.DemandExamine = row["demandExamine"].(uint8)
		warehouse.DeptId = row["deptId"].(uint)
		warehouse.IsOversea = row["isOversea"].(uint8)
		warehouse.SpecialDataExamine = row["specialDataExamine"].(uint8)
		warehouse.SpecialDemandExamine = row["specialDemandExamine"].(uint8)
		warehouse.StockYunCode = row["stockYunCode"].(string)
		warehouse.Type = row["Type"].(uint8)
		warehouse.WarehouseType = row["warehouseType"].(uint8)
		err = db.FirstOrCreate(&warehouse).Error
		if err != nil {
			panic(err)
		}
	})

	warehouseDept := WarehouseDept{}
	err = db.Last(&warehouseDept).Error
	ifs.FetchWarehouseDept(mysqlDb, warehouseDept.ID, func(row map[string]interface{}) {
		warehouseDept := WarehouseDept{}
		warehouseDept.Name = row["name"].(string)
		warehouseDept.State = row["status"].(uint8)
		warehouseDept.Code = row["sn"].(string)
		warehouseDept.ID = row["id"].(uint)
		err = db.FirstOrCreate(&warehouseDept).Error
		if err != nil {
			panic(err)
		}
	})

	category := Category{}
	err = db.Last(&category).Error
	ifs.FetchCategory(mysqlDb, category.ID, func(row map[string]interface{}) {
		category := Category{}
		category.ID = uint(row["id"].(int64))
		category.State = uint8(row["is_enable"].(int64))
		category.Name = fmt.Sprintf(`%s`, row["name_ch"])
		category.Path = fmt.Sprintf(`%s`, row["path"])
		category.CatalogType = fmt.Sprintf(`%s`, row["catalog_type"])
		category.DisplayAlone = uint8(row["display_alone"].(int64))
		category.EditSamePriceReadonly = uint8(row["edit_same_price_readonly"].(int64))
		category.IsEditSamePrice = uint8(row["is_edit_same_price"].(int64))
		category.IsLimitSyn = uint8(row["is_limit_syn"].(int64))
		category.IsMaking = uint8(row["is_making"].(int64))
		category.IsProviderSee = uint8(row["is_provider_see"].(int64))
		category.Level = uint8(row["level"].(int64))
		category.NeedQc = uint8(row["need_qc"].(int64))
		category.OrderID = uint(row["order_id"].(int64))
		category.ParentID = uint(row["parent_id"].(int64))
		category.PictureURL = fmt.Sprintf(`%s`, row["picture_url"])
		category.SampleAlready = uint8(row["sample_already"].(int64))
		category.UpdateUser = fmt.Sprintf(`%s`, row["update_user"])
		category.WhCode = fmt.Sprintf(`%s`, row["wh_code"])
		err = db.FirstOrCreate(&category).Error
		if err != nil {
			log.Fatalln(err)
		}
	})
}
