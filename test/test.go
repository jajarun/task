package main

import (
	"fmt"
	"task/models"
)

func main() {
	//user := models.User{}
	//models.FindByPrimaryKey(&user, "18")
	//fmt.Println(user)
	//user.UserName = "jajatest"
	//models.Save(&user)
	////fmt.Println(user)

	users := []models.User{}
	models.Find(&users)
	fmt.Println(users)

}
