/*
* @Author: wuhialin
* @Date:   2018-01-16 10:29:39
* @Last Modified by:   wuhailin
* @Last Modified time: 2018-01-22 14:35:00
 */
package model

import (
	"github.com/jinzhu/gorm"
)

type Model struct {
	gorm.Model
	State uint8 `gorm:"NOT NULL;default:1"`
}
