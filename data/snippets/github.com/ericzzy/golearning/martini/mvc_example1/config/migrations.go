package main

import (
	"fmt"
	"os"

	"../models"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func main() {
	db, err := gorm.Open("mysql", "root:passw0rd@/martini_mvc?charset=utf8&parseTime=True&loc=Local")
	// Create table
	if err != nil {
		fmt.Printf("Faild to connect to database with error: %v \n", err)
		os.Exit(1)
	}
	db.CreateTable(&models.Activity{})
}
