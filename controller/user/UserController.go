package user

import (
	"fmt"
	"task/components"
	"task/controller/base"
	"task/models"
)

func Login(ch *base.ControllerHandle) {
	bodyMap := ch.GetMapBody()
	mobile := bodyMap["mobile"].(string)
	password := bodyMap["password"].(string)
	fmt.Println(mobile)
	fmt.Println(password)
	user := models.FindUserByMobile(mobile)
	if user.ID <= 0 {
		ch.ReturnError(1, "用户不存在")
		return
	}
	if !components.Login(password, user.Password) {
		ch.ReturnError(1, "密码有误")
		return
	}
	ch.ReturnData()
}

func Register(ch *base.ControllerHandle) {
	bodyMap := ch.GetMapBody()
	userName := bodyMap["user_name"].(string)
	mobile := bodyMap["mobile"].(string)
	password := bodyMap["password"].(string)
	user := models.User{
		UserName: userName,
		Password: components.EncryptPassword(password),
		Mobile:   mobile,
	}
	models.Create(&user)
	ch.ReturnData(user)
}

func Info(ch *base.ControllerHandle) {
	id := ch.GetQuery("id")
	if id == "" {
		ch.ReturnError(base.ERROR_DEFAULT, "errorParam")
		return
	}
	user := models.User{}
	models.FindByPrimaryKey(&user, id)

	ch.ReturnData(user)
}
