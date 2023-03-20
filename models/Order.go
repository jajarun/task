package models

type Order struct {
	ModelBase

	Id     uint `gorm:"primarykey"`
	UserId uint `gorm:"primarykey"`
}
