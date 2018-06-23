package main

import (
	"fmt"
	"github.com/gophergala/echodb/db"
	"github.com/gophergala/echodb/dbhttp"
	"time"
)

func temporaryTodoCleaner(echodb *db.Database) {
	for _ = range time.Tick(10 * time.Minute) {
		echodb.Delete("todo")
		echodb.Create("todo")
	}
}

func main() {
	echodb, err := db.OpenDatabase("/tmp/echodb")
	if err != nil {
		fmt.Printf("Error: %v", err)
		panic("exit")
	}

	echodb.Create("books")
	echodb.Create("todo")

	go temporaryTodoCleaner(echodb)

	dbhttp.Start()
}
