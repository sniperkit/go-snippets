package model

type HtmlSelector struct {
	Model

	Selector string `gorm:"not null;default:''"`
	Sorting  uint8  `gorm:"not null;default:0"`
	Eq       int8   `gorm:"not null;default:-1"`

	Source   Source
	SourceId uint `gorm:"index;not null;default:0"`

	HtmlSelectorGroup HtmlSelectorGroup `gorm:"ForeignKey:ProfileRefer"`
	GroupId           uint              `gorm:"index;not null;default:0"`
}

type HtmlSelectorGroup struct {
	Model

	Name      string `gorm:"not null;default:'';size:50;"`
	TableName string `gorm:"not null;default:'';size:50"`

	Source   Source
	SourceId uint `gorm:"not null;default:0;index"`

	HtmlSelectors []HtmlSelector
}
