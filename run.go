package main

import (
	"github.com/spf13/viper"
	"net/http"
	"task/route"
)

func initConfig() {
	viper.SetConfigFile("config/app.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic("load config err" + err.Error())
	}
}

func main() {
	initConfig()

	route.WebRouteInit()
	route.AccRouteInit()

	_ = http.ListenAndServe(":8181", nil)
}
