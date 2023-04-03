package main

import (
	"github.com/spf13/viper"
	"net/http"
	"task/route"
	"task/ws"
)

func main() {

	//fmt.Println("app start")
	viper.SetConfigFile("config/app.yaml")
	_ = viper.ReadInConfig()
	////fmt.Println(viper.Get("mysql.username"))
	//
	ws.Run()
	route.WebRouteInit()
	//route.AccRouteInit()

	_ = http.ListenAndServe(":8181", nil)
}
