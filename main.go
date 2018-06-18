package main

import (
	"fmt"
	"os"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/uphy/gorm-examples/checkExistence"
	"github.com/uphy/gorm-examples/compositePrimaryKey"
	"github.com/uphy/gorm-examples/hasMany"

	"github.com/jinzhu/gorm"
)

type Example func(db *gorm.DB)

var examples = map[string]Example{
	"Has many":                  hasMany.Example,
	"Check if the record exist": checkExistence.Example,
	"Composite Primary Key":     compositePrimaryKey.Example,
}

func main() {
	for key, _ := range examples {
		fmt.Println("-------------------------------------")
		fmt.Println(key)
		fmt.Println("-------------------------------------")
		run(key)
	}
}

func run(name string) {
	ex := examples[name]
	file := "test.db"
	db, err := gorm.Open("sqlite3", file)
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()
	defer os.Remove(file)
	ex(db)
	if db.Error != nil {
		fmt.Printf("Error: %s\n", db.Error)
	}
}
