package main

import (
	"fmt"
	"task/tasks"
	"time"
)

func main() {
	//user := models.User{}
	////
	////models.FindFirst(&user, 2)
	////fmt.Printf("data:%v \n", user.UserName)
	//models.FindByPrimaryKey(&user, "2")
	//
	//user.Password = "jaja1234567"
	//models.Save(&user)
	////fmt.Println("user data:", user)
	//models.FindByPrimaryKey(&user, "2")

	//user.Password = "45678"
	//order := models.Order{}
	//models.AutoMigrate(&order)

	//for i := 1; i <= 5; i++ {
	//	go func(index int) {
	//		dbClient := db.GetInstanceDb()
	//		user := models.User{}
	//		dbClient.First(&user, index)
	//		fmt.Printf("index %d name %s \n", index, user.UserName)
	//	}(i)
	//}
	//time.Sleep(time.Second * 5)

	tasks.HandleQueue("go:test:queue:task", func(item interface{}) {
		fmt.Println("start do item:", item)
		time.Sleep(3 * time.Second)
		fmt.Println("end do item:", item)
	})

}
