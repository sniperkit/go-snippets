/*
* @Author: wuhialin
* @Date:   2018-01-16 10:47:29
* @Last Modified by:   wuhialin
* @Last Modified time: 2018-01-16 11:11:14
 */
package model

type WarehouseDept struct {
	Model
	Name         string  `gorm:"not null; default:'';size:50"`
	Code         string  `gorm:"not null;default:'';size:7;unique_index"`
	MinOrderRate float32 `gorm:"not null;default:0"`
}
