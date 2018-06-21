package model

type Source struct {
	Model
	Name   string `gorm:"size:20;index;not null;default:''"`
	Domain string `gorm:"index;not null;default:''"`
	Home   string `gorm:"not null;default:''"`
}
