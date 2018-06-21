package models

import (
	fgorm "github.com/henrylee2cn/faygo/ext/db/gorm"
	"github.com/jinzhu/gorm"
)

const (
	AppName = "app"
)

var defaultDb *gorm.DB

func init() {
	db, ok := fgorm.DB()
	if !ok {
		panic("load db failure")
	}
	db.SingularTable(true)
	defaultDb = db
}

type Model struct {
	//
}

func (t *Model) DB() *gorm.DB {
	return defaultDb.Model(t)
}
