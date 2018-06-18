package compositePrimaryKey

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	Name  string `gorm:"primary_key"`
	Email string `gorm:"primary_key"`
}

type MutedLogger struct {
}

func (m MutedLogger) Print(v ...interface{}) {
	// do nothing
}

func Example(db *gorm.DB) {
	db.AutoMigrate(&User{})
	user1 := User{
		Name:  "a",
		Email: "a@foo.com",
	}
	if err := db.Create(&user1).Error; err != nil {
		panic(err)
	}

	user2 := User{
		Name:  "a",
		Email: "a@foo.com",
	}
	db.SetLogger(MutedLogger{})
	if err := db.Create(&user2).Error; err != nil {
		// expected behavior
	} else {
		panic("duplicated record must not be inserted")
	}
	db.SetLogger(gorm.Logger{})

	user3 := User{
		Name:  "a",
		Email: "a@bar.com",
	}
	if err := db.Create(&user3).Error; err != nil {
		panic(err)
	} else {
		// expected behavior
	}

	allUsers := []User{}
	db.Find(&allUsers)
	fmt.Println("Created users:")
	for _, v := range allUsers {
		fmt.Printf("%#v\n", v)
	}
}
