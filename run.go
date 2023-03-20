package main

import (
	"task/models"
)

func main() {
	user := models.User{}
	//
	//models.FindFirst(&user, 2)
	//fmt.Printf("data:%v \n", user.UserName)
	models.FindByPrimaryKey(&user, "2")

	user.Password = "jaja1234567"
	models.Save(&user)
	//fmt.Println("user data:", user)
	models.FindByPrimaryKey(&user, "2")

	//user.Password = "45678"
}
