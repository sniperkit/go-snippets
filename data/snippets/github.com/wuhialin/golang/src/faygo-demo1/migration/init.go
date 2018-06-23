package migration

import (
	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/gorm"
	sourceOrm "github.com/jinzhu/gorm"
)

var db *sourceOrm.DB

func Run() {
	var ok bool
	db, ok = gorm.DB("ifs")
	if !ok {
		faygo.Error("load ifs db failure")
		return
	}
	if err := db.DB().Ping(); err != nil {
		faygo.Error(err)
		return
	}

	go warehouse()
}
