package models

type User struct {
	ModelBase

	ID         uint   `gorm:"primarykey" json:"id"`
	UserName   string `json:"userName"`
	Mobile     string `json:"mobile"`
	Password   string `json:"password"`
	CreateTime uint   `gorm:"autoCreateTime" json:"createTime"`
	UpdateTime uint   `gorm:"autoUpdateTime" json:"updateTime"`
}

//func (u *User) AfterSave(tx *gorm.DB) (err error) {
//	FlushCache(u)
//	return nil
//}

//func (u *User) AfterQuery(tx *gorm.DB) (err error) {
//	fmt.Println(u)
//	return nil
//}

//func (u *User) getPrimaryKey() string {
//	return "ID"
//}

func (u *User) IsModelCache() bool {
	return false
}

func (u *User) GetRevisionClue() string {
	return "user_id"
}

func FindUserByMobile(mobile string) User {
	user := User{}
	db := GetInstanceDb()
	db.Where("mobile = ?", mobile).First(&user)
	return user
}
