/*
* @Author: wuhailin
* @Date:   2018-01-20 14:53:17
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-20 14:53:39
 */
package model

type Category struct {
	Model
	Name                  string `gorm:"not null;default:'';size:100"`
	ParentID              uint   `gorm:"not null;default:0;index"`
	Level                 uint8  `gorm:"not null;default:0"`
	Path                  string `gorm:"not null;default:'';size:50;index"`
	OrderID               uint   `gorm:"not null;default:0"`
	DisplayAlone          uint8  `gorm:"not null;default:0"`
	SampleAlready         uint8  `gorm:"not null;default:0"`
	NeedQc                uint8  `gorm:"not null;default:0"`
	IsEditSamePrice       uint8  `gorm:"not null;default:0"`
	EditSamePriceReadonly uint8  `gorm:"not null;default:0"`
	IsMaking              uint8  `gorm:"not null;default:0"`
	CatalogType           string `gorm:"not null;default:'';size:30"`
	PictureURL            string `gorm:"not null;default:'';size:200"`
	UpdateUser            string `gorm:"not null;default:'';size:30"`
	IsLimitSyn            uint8  `gorm:"not null;default:0"`
	IsProviderSee         uint8  `gorm:"not null;default:0"`
	WhCode                string `gorm:"not null;default:'';size:45"`
}
