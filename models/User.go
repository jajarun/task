package models

type User struct {
	ModelBase

	ID         uint `gorm:"primarykey"`
	UserName   string
	Mobile     string
	Password   string
	CreateTime uint `gorm:"autoCreateTime"`
	UpdateTime uint `gorm:"autoUpdateTime"`
}

func (u User) isModelCache() bool {
	return false
}

func FindUserByMobile(mobile string) User {
	user := User{}
	db := GetInstanceDb()
	db.Where("mobile = ?", mobile).First(&user)
	return user
}

//
//func (u User) getRevisionClue() string {
//	return "user_id"
//}
