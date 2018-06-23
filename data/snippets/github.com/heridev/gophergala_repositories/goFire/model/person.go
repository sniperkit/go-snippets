package model

import "time"

type Person struct {
	PersonID     	int  		`db:"PersonID"`
	CompanyID  		int 		`db:"CompanyID"`
	FirstName 		string		`db:"FirstName"`
	LastName 		string 		`db:"LastName"`
	Email 			string 		`db:"Email"`
	CreateDate   	time.Time 	`db:"CreateDate"`
	CreatedPersonID	int 	    `db:"CreatedPersonID"`
	EditedDate   	time.Time 	`db:"EditedDate"`
	EditedPersonID  int    		`db:"EditedPersonID"`
}
