package main

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"log"
)

func main() {
	consulClient, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		log.Fatal(err)
	}
	//services, err := consulClient.Agent().Services()
	//services, err := consulClient.Agent().ServicesWithFilter("Service == MyService")
	service, meta, serr := consulClient.Health().Service("MyService", "", true, nil)
	log.Printf("结束svr%+v,meta:%+v, err:%v\n", service, meta, serr)
	for _, s := range service {
		fmt.Printf("服务：%+v\n", s.Service)
	}
	//if err != nil {
	//	log.Fatal(err)
	//}
	//for _, v := range services {
	//	fmt.Println("Got Service:", v)
	//	fmt.Println("Got Service:", v.Service)
	//}
	//fmt.Printf("Got Service: ip-%s, port-%d", service.Address, service.Port)
	//fmt.Println("Got Service:", services)
}
