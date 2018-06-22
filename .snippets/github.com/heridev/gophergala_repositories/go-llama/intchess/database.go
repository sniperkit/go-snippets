package intchess

import (
	"code.google.com/p/go.crypto/bcrypt"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var dbGorm *gorm.DB

func ConnectToDatabase() error {
	dbGormConnection, err := gorm.Open("mysql", dbConnString) //NOT INCLUDED IN REPO
	if err != nil {
		return err
	}
	dbGorm = &dbGormConnection
	return nil
}

func CreateDatabaseTables() {
	fmt.Printf("I am attempting to create the database tables!")
	fmt.Printf("Dropping (if exists) and creating Users table...\n")
	dbGorm.DropTableIfExists(&User{})
	dbGorm.CreateTable(&User{})

	fmt.Printf("Dropping (if exists) and creating ChessGame table...\n")
	dbGorm.DropTableIfExists(&ChessGame{})
	dbGorm.CreateTable(&ChessGame{})

	dbGorm.DropTableIfExists(&UserRankChange{})
	dbGorm.CreateTable(&UserRankChange{})

	dbGorm.DropTableIfExists(&GameMove{})
	dbGorm.CreateTable(&GameMove{})

	//create me a default user
	fmt.Printf("Adding default test users to database...\n")
	pass, _ := bcrypt.GenerateFromPassword([]byte("test"), 3)
	u := User{
		Username:    "test",
		AccessToken: string(pass),
		IsAi:        false,
		VersesAi:    true,
	}
	dbGorm.Create(&u)
	u = User{
		Username:    "test2",
		AccessToken: string(pass),
		IsAi:        false,
		VersesAi:    true,
	}
	dbGorm.Create(&u)
	fmt.Printf("Done!")
}
