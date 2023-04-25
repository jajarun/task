package main

import (
	"fmt"
	"github.com/spf13/viper"
	"task/models"
)

func initConfig() {
	viper.SetConfigFile("/Users/majh/GolandProjects/task/config/app.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("load config err" + err.Error())
	}
}

func main() {
	initConfig()
	//users := []models.User{}
	db := models.GetInstanceDb()
	//db.First(&users)
	//fmt.Println("result:", users)

	user := models.User{}
	db.First(&user)
	fmt.Println("result:", user)

}
