package main

import (
	"strconv"
	"task/models"
	"task/tasks"
)

func main() {
	//user := models.User{
	//	UserName: "李四",
	//	Password: "2423412",
	//}
	//models.Create(&user)
	//
	//redis := redisClient.GetInstanceRedis()
	//redis.LPush("test:queue", 6)
	//redis.LPush("test:queue", 7)

	tasks.HandleQueue("test:queue", func(item interface{}) {
		var id int
		switch item.(type) {
		case int:
			id = item.(int)
			break
		case string:
			id, _ = strconv.Atoi(item.(string))
			break
		default:
			panic("err item")
		}
		user := models.User{}
		models.FindFirst(&user, id)
		models.FindFirst(&user, id)
	})
}
