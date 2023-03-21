package models

import "gorm.io/gorm"

type Order struct {
	gorm.Model
	Id     uint `gorm:"primarykey"`
	UserId uint `gorm:"primarykey"`
}
