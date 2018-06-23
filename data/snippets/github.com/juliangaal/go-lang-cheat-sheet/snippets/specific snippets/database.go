package main

import (
	"log"

	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var recipeDB *gorm.DB

//Struct to define entries into Database
type Data struct {
	ID     int `gorm:"primary_key"`
	Name   string
	Gender string
}

func connectToDatabase() {
	var err error

	//parameters may be adjusted, add passwords, ssl etc
	exampleDB, err = gorm.Open("postgres", "user=Julian dbname=examplename sslmode=disable")

	if err != nil {
		log.Printf("Cannot connect to database \n")
		return
	}

}

func feedData() {
	var (
		data  []Data
		entry = Data{Name: "Julian", Gender: "apache"}
	)

	exampleDB.AutoMigrate(&Data{})

	//Checks if entry already exists in database
	err := exampleDB.Where("exampleparam = ?", exampledata).Find(&data).Error
	if err != nil {
		log.Printf("Couldn't Connect to database when adding entry: %v", err)
	}

	//If returned 'data' is not empty, entry already exists
	if len(data) > 0 {
		log.Printf("entry already exists. not added %v %v", recipe[0], recipe[1])
	} else {
		exampleDB.Create(&entry)
	}
}

func main() {
	connectToDatabase()
	feedData()

	defer recipeDB.Close()
}
