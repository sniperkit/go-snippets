package database

import (
  "fmt"
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
  "os"
)

var db *gorm.DB = nil

func DB() *gorm.DB {
  if db == nil {
    connection, err := gorm.Open("postgres", os.Getenv("GOTOOLBOX_POSTGRES_URL"))
    if err != nil {
      panic(fmt.Sprintf("Could not connect to database"))
    }
    connection.LogMode(true)
    db = &connection
  }
  return db
}
