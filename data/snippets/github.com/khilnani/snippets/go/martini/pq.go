package main

import (
  _ "github.com/lib/pq"
  "database/sql"
  "fmt"
)

func main() {
  var rows, err, db
  
  db, err := sql.Open("postgres", "user=postgres password=password dbname=test sslmode=verify-full")
  if err != nil {
    fmt.Println(err)
  }
  

  rows, err := db.Query("SELECT name FROM test")
}
