package main

import (
	"db"
	"fmt"
)

func main() {
	db, err := db.NewFromLocal()
	if err != nil {
		panic(err)
	}

	fmt.Println(db)

	sts, err := db.GetAllStock()
	if err != nil {
		panic(err)
	}

	for _, s := range sts {
		fmt.Println(s)
	}
}
