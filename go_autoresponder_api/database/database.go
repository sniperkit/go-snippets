package database

import (
  "github.com/jinzhu/gorm"
  _ "github.com/lib/pq"
  "github.com/heridev/go_autoresponder_api/utils"
  "github.com/heridev/go_autoresponder_api/models"
)

var DbInstance gorm.DB

func InitDb() {
  // connect to db using standard Go database/sql API
  var err error
  DbInstance, err = gorm.Open("postgres", "dbname=hmail sslmode=disable")
  utils.PanicIf(err)

  DbInstance.AutoMigrate(&models.Subscriber{},
                         &models.EmailList{},
                         &models.Autoresponder{})

  DbInstance.DB().SetMaxIdleConns(20)

  autoresponder := models.Autoresponder{
    Title:            "title autoresponder",
    Description:      "description goes here",
    Lists:          []models.EmailList{
                      {
                        Title: "template 1",
                        Content: "<h1>Template 1</h1>",
                      },
                    },
  }

  DbInstance.Create(&autoresponder)
}
