package model

import (
	"strings"
	"github.com/jinzhu/gorm"
	"net/url"
	"log"
)

type Warehouse struct {
	Name string `gorm:"size:50;not null;default:''"`
	Code string `gorm:"size:50;not null;default:'';unique_index"`

	UserModel
}

func (t *Warehouse) SearchList(v url.Values) *gorm.DB {
	curDb := db.Model(t)
	if strings.TrimSpace(v.Get("code")) != "" {
		curDb = curDb.Where(Warehouse{Code: v.Get("code")})
	}
	log.Println(v.Get("name"))
	if strings.TrimSpace(v.Get("name")) != "" {
		curDb = curDb.Where("name LIKE ?", v.Get("name")+"%")
	}
	return curDb
}
