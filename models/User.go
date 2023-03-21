package models

import (
	"fmt"
)

type User struct {
	ModelBase

	ID         uint `gorm:"primarykey"`
	UserName   string
	Password   string
	CreateTime uint `gorm:"autoCreateTime"`
	UpdateTime uint `gorm:"autoUpdateTime"`
}

func (u User) getPrimaryId() string {
	return fmt.Sprintf("%d", u.ID)
	//return map[string]string{
	//	"ID": fmt.Sprintf("%d", u.ID),
	//}
}

//func (u User) isModelCache() bool {
//	return true
//}
//
//func (u User) getRevisionClue() string {
//	return "user_id"
//}
