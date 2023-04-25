package main

import (
	"fmt"
	"github.com/micro/micro/v3/service"
	serviceMicro "task/micro-proto"
)

func main() {
	srv := service.New(func(o *service.Options) {
		o.Name = "hello_server"
		o.Version = "v1.0.0"
	})
	//服务初始化
	//srv.Init()
	//注册
	err := serviceMicro.RegisterGreeterHandler(srv.Server(), new(serviceMicro.HelloworldService))
	if err != nil {
		fmt.Println("register err:", err)
		return
	}
	//
	////运行
	err = srv.Run()
	if err != nil {
		fmt.Println("run err:", err)
	}
}
