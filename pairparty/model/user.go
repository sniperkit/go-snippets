package model

import (
	"time"
)

type User struct {
	Id           int64
	Username     string `sql:"not null"`
	Name         string `sql:"size:255"`
	Email        string `sql:"not null"`
	Password     string `sql:"not null"` //Yeah this is lame
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    time.Time
}
