package main

import (
	"fmt"
	"querybuilder/db"
)

func main() {

	//create a db object // this will remain same for the rest of the application
	db1, err1 := db.New(db.MYSQL) //db.SQLITE // db.MONGODB
	if err1 != nil {
		fmt.Println("Error1 in db1: " + err1.Error())
		return
	}

	//set a table name
	db1.Table("users")

	users1 := db1.Get() //fetch all the users from users table
	fmt.Println(users1)

	user1 := db1.Find(101) // fetch a user with id 101
	fmt.Println(user1)

	removeUser101 := db1.Query("delete where id=101")
	fmt.Println(removeUser101)

	//repeat the same code just rename the variables
	//This time db.New will return the old object (as it's following singleton)
	db2, err2 := db.New(db.MYSQL)
	if err2 != nil {
		fmt.Println("Error2 in db2: " + err2.Error())
		return
	}

	users2 := db2.Get() //this time we did not tell the table name// but it'll use the previous table name as it refer to the old object
	fmt.Println(users2)
}
