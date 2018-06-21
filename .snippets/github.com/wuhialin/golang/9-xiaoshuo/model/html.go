package model

type Html struct {
	Model

	Url   Url
	UrlId uint `gorm:"not null;default:0;unique_index"`

	Data string `gorm:"not null"`
}

type HtmlParse struct {
	Model

	HtmlId uint `gorm:"not null;default:0;index"`

	GroupId uint `gorm:"not null;default:0;index"`
}
