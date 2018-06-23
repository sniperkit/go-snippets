package model

type Group struct {
	GroupID     int    `db:"GroupID"`
	GroupTitle  string `db:"GroupTitle"`
	Description string `db:"Description"`
}
