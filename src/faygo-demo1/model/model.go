package model

import (
	"github.com/jinzhu/gorm"
	gorm2 "github.com/henrylee2cn/faygo/ext/db/gorm"
)

const AppName = "app"

var db *gorm.DB

func init() {
	var ok bool
	if db, ok = gorm2.DB(); !ok {
		panic("load db failure")
	}
	if err := db.DB().Ping(); err != nil {
		panic(err)
	}
	db.SingularTable(true)

	db.AutoMigrate(new(Warehouse), new(User))

	initUser()
}

type Model struct {
	gorm.Model
	State uint8 `gorm:"not null;default:1"`
}

type UserModel struct {
	Model

	CreateUser uint `gorm:"not null"`
	UpdateUser uint `gorm:"not null"`
}

func (t *Model) Save(m interface{}) *gorm.DB {
	return db.Save(m)
}

func (t *Model) Last(m interface{}) *gorm.DB {
	return db.Last(m)
}

func (t *Model) DB() *gorm.DB {
	return db
}

func initUser() {
	count := 0
	user := new(User)
	if err := db.Model(user).Where("id = ?", 1).Count(&count).Error; err != nil {
		panic(err)
	}
	if count == 0 {
		user.ID = 1
		user.Username = "admin"
		user.Salt = user.GenerateRandom()
		user.Password = user.GeneratePassword("admin")
		if err := db.Save(user).Error; err != nil {
			panic(err)
		}
	}
}
