/*
* @Author: wuhialin
* @Date:   2018-01-16 11:11:58
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-25 10:34:14
 */
package model

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

const (
	AppName = "app"
)

var db *gorm.DB
var mysqlDb *gorm.DB

func init() {
	var err error
	db, err = gorm.Open("postgres", "host=127.0.0.1 user=postgres dbname=postgres sslmode=disable password=postgres")
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	err = db.DB().Ping()
	if err != nil {
		panic(err)
	}
	db.SingularTable(true)
	db.AutoMigrate(new(Warehouse), new(WarehouseDept), new(Category))

	mysqlDb, err = gorm.Open("mysql", "ifs_ro:dYDIiSFY@tcp(192.168.7.165)/ifs?charset=utf8")
	if err != nil {
		panic(err)
	}
	err = mysqlDb.DB().Ping()
	if err != nil {
		panic(err)
	}

	migrateIfs()
}

func DB() *gorm.DB {
	return db
}
