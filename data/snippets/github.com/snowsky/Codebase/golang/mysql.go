package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "ocp_read:trejac9enAxa@tcp(cam-vm-db1a.tpn.thinkingphones.net:3306)/ocp")
	if err != nil {
		log.Fatal("Failed to connect database!")
	}
	defer db.Close()

	rows, err := db.Query("select code from data_centers")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for rows.Next() {
		var code string
		if err := rows.Scan(&code); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Done!" + code)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}
