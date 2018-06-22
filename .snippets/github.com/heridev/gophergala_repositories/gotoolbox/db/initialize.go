package main

import (
  "fmt"
  "github.com/gophergala/gotoolbox/models"
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
  "os"
)

func main() {
  fmt.Println("Seeding Database")
  db, _ := gorm.Open("postgres", os.Getenv("GOTOOLBOX_POSTGRES_URL"))
  db.LogMode(true)

  db.DB()

  db.DropTableIfExists(&models.User{})
  db.DropTableIfExists(&models.Project{})
  db.DropTableIfExists(&models.Category{})

  // create database
  db.AutoMigrate(&models.User{})
  db.AutoMigrate(&models.Project{})
  db.AutoMigrate(&models.Category{})

  // some nice seeding for the categories
  router := models.Category{Name: "Router", Description: "net/http router"}
  db.Save(&router)

  webframeworks := models.Category{Name: "Web Frameworks", Description: "Webframeworks available"}
  db.Save(&webframeworks)

  json := models.Category{Name: "JSON", Description: "JSON Parsing"}
  db.Save(&json)

  kvs := models.Category{Name: "Key Value Stores", Description: "Key Value Stores"}
  db.Save(&kvs)

  orm := models.Category{Name: "ORM", Description: "Object relational mapping"}
  db.Save(&orm)

  te := models.Category{Name: "Template Engines", Description: "Template Engines"}
  db.Save(&te)

  ff := models.Category{Name: "File Formats", Description: "File Formats"}
  db.Save(&ff)

  cx := models.Category{Name: "Compression", Description: "Compression"}
  db.Save(&cx)

  oauth := models.Category{Name: "OAuth", Description: "OAuth"}
  db.Save(&oauth)

  crypto := models.Category{Name: "Crypto", Description: "Crypto"}
  db.Save(&crypto)
}
