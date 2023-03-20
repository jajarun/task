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

func (u User) GetPrimaryId() string {
	return fmt.Sprintf("%d", u.ID)
}

func (u User) IsModelCache() bool {
	return true
}

func (u User) GetRevisionClue() string {
	return "user_id"
}
