// models/basemodel.go

package models

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type BaseModel struct {
	ID        uint `gorm:"primary_key" form:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
