/*
* @Author: wuhailin
* @Date:   2018-01-18 14:46:23
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-18 14:59:22
 */
package model

import "github.com/henrylee2cn/faygo"

type Warehouse struct {
	Model
	Name                 string `gorm:"not null;default:'';size:50"`
	Code                 string `gorm:"not null;default:'';size:50;unique_index"`
	AutoSend             uint8  `gorm:"not null;default:0"`
	Type                 uint8  `gorm:"not null;default:0"`
	WarehouseType        uint8  `gorm:"not null;default:0"`
	IsOversea            uint8  `gorm:"not null;default:0"`
	DeptId               uint   `gorm:"not null;default:0"`
	SpecialDataExamine   uint8  `gorm:"not null;default:0"`
	SpecialDemandExamine uint8  `gorm:"not null;default:0"`
	DemandExamine        uint8  `gorm:"not null;default:0"`
	StockYunCode         string `gorm:"not null;default:'';size:50"`
	DemandAutoPass       uint8  `gorm:"not null;default:0"`
}

