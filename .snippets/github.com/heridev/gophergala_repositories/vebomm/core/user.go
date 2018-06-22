package core

import (
	"time"
)

type User struct {
	Id        int64  `validate:"-"`
	Mmr       int    `validate:"-"`
	Username  string `validate:"min=2,regexp=^[a-zA-Z_]*$"`
	Title     string `validate:"regexp=^[a-zA-Z_.]*$"`
	Password  string `validate:"nonzero,min=4,regexp=^.+$"`
	Email     string `validate:"nonzero,regexp=^[a-zA-Z_.]+@[a-zA-Z_.]+$"`
	BirthYear int    `validate:"nonzero,min=1900,max=2100"`
	Gender    int    `validate:"min=0,max=1"`
}

type RegisterResult struct {
	ValOk       bool
	DupUsername bool
	DupEmail    bool
}

type LoginResult struct {
	Ok   bool
	User *User
}

func (u User) Age() int {
	return time.Now().Year() - u.BirthYear
}
