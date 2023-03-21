package main

import "task/models"

func main() {
	user := models.User{}
	models.FindByPrimaryKey(&user, "6")
	user.UserName = "张三1号"
	models.Save(&user)
}
