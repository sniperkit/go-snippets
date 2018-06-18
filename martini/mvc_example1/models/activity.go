package models

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jinzhu/gorm"
)

type Activity struct {
	BaseModel
	Title       string `form:"title" binding:"required"`
	Description string `form:"description"`
	City        string `form:"city" binding:"required"`
	Address     string `form:"address"`
	Contract    string `form:"contract"`
}
