package hasMany

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type User struct {
	ID    uint `gorm:"primary_key"`
	Name  string
	Roles []Role
}

type Role struct {
	Name   string
	UserID uint
}

func Example(db *gorm.DB) {
	// init db
	db.AutoMigrate(&User{})
	db.AutoMigrate(&Role{})
	db.Model(&User{}).Related(&Role{})
	// create user1 in 'users' table and role1, role2 in 'roles' table.
	user := User{
		Name: "user1",
		Roles: []Role{
			Role{Name: "role1"},
			Role{Name: "role2"},
		},
	}
	db.Create(&user)

	// find the created user.  for eager loading, 'Preload("Roles")' is required.
	found := &User{}
	db.Preload("Roles").Find(&found)
	fmt.Println(found)
}
