package conf

import (
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"os"
)

var DB *gorm.DB

func SetupDB() *gorm.DB {
	connectionString := os.Getenv("DATABASE_URL")

	if connectionString == "" {
		connectionString = "dbname=gopherstalker sslmode=disable"
	}

	db, err := gorm.Open("postgres", connectionString)
	db.LogMode(true)
	PanicIf(err)
	DB = &db
	return DB
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
