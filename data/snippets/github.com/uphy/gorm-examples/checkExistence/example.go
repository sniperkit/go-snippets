package checkExistence

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID   uint `gorm:"primary_key"`
	Name string
}

func Example(db *gorm.DB) {
	db.AutoMigrate(&User{})
	user := User{
		ID:   1,
		Name: "foo",
	}
	db.Create(&user)

	// find existing 'foo'
	found := User{}
	if result := db.Find(&found, 1); result.RecordNotFound() {
		panic("record should be found")
	}
	fmt.Printf("record found: %#v\n", found)

	if result := db.Find(&found, 1234); result.RecordNotFound() == false {
		panic("record should not be found")
	}
	fmt.Println("record not found")
}
