package models

type Employee struct {
	ID         uint `gorm:"primary_key"`
	Username   string
	Mobile     string
	Name       string
	Tel        string
	Mobile     string
	Email      string
	CreateTime uint
	LastTime   uint
	LastIp     string
	JobName    string

	Model
}
