package model

type GroupDetail struct {
	GroupDetailID int `db:"GroupDetailID"`
	PersonID	  int `db:"PersonID"`
	GroupID 	  int `db:"GroupID"`
}
