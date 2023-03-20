package models

import (
	"fmt"
	"gorm.io/gorm"
)

type User struct {
	ModelBase

	gorm.Model
	UserName string
	Password string
}

func (u User) GetPrimaryId() map[string]string {
	//return fmt.Sprintf("%d", u.ID)
	return map[string]string{
		"ID": fmt.Sprintf("%d", u.ID),
	}
}

func (u User) IsModelCache() bool {
	return true
}

func (u User) GetRevisionClue() string {
	return "user_id"
}
