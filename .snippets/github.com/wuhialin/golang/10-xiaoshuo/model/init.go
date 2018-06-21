package model

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"os"
	"path/filepath"
	"time"
)

var db *gorm.DB

func init() {
	var err error
	source := "host=localhost user=postgres dbname=postgres sslmode=disable password=123456"
	db, err = gorm.Open("postgres", source)
	if err != nil {
		panic(err)
	}

	db.SingularTable(true)

	db.DB().SetMaxOpenConns(1000)
	db.DB().SetMaxIdleConns(10)

	initAutoMigrate()
	initAutoData()
}

func initAutoMigrate() {
	db.AutoMigrate(&Source{}, &Book{}, &BookItem{}, &BookPage{})
}

func initAutoData() {
	curPath, _ := os.Getwd()
	logPath := filepath.Join(curPath, "log")
	f, err := os.OpenFile(filepath.Join(logPath, "data.lock"), os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fInfo, _ := f.Stat()
	if fInfo.Size() > 0 {
		return
	}
	tmpDb := db.Begin()
	source := Source{
		Name:   "61读书网",
		Domain: "m.61xsw.com",
		Home:   "http://m.61xsw.com",
	}
	tmpDb.Save(&source)
	tmpDb.Commit()
	f.WriteString(fmt.Sprintln(time.Now().Unix()))
}

func DB() *gorm.DB {
	return db.Scopes(stateWhere).Debug()
}

type Model struct {
	gorm.Model
	State uint8 `gorm:"not null;default:1"`
}

func stateWhere(db *gorm.DB) *gorm.DB {
	return db.Where("state = ?", 1)
}
