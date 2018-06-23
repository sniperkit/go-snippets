package model

import "time"

type User struct {
	UserID      int32
	FirstName   string `sql:"type:varchar(75);"`
	LastName    string `sql:"type:varchar(75);"`
	Email       string `sql:"type:varchar(150);"`
	Login       string `sql:"type:varchar(150);"`
	AccessToken string `sql:"type:varchar(4000);"`
	DateCreated time.Time
	Visible     bool
}
