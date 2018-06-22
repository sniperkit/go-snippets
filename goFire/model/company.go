package model

import "time"

type Company struct {
	CompanyID    	int       	`db:"CompanyID"`
	CompanyName 	string    	`db:"CompanyName"`
	CreatedDate   	time.Time 	`db:"CreatedDate"`
	CreatedPersonID int			`db:"CreatedPersonID"`
	EditedDate   	time.Time 	`db:"EditedDate"`
	EditedPersonID 	int			`db:"EditedPersonID"`
}
