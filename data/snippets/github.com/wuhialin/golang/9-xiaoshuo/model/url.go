package model

type Url struct {
	Model
	Url string `gorm:"size:255;not null;default:''"`

	Source   Source
	SourceId uint `gorm:"not null;default:0;index"`
}
