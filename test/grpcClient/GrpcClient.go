package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	service "task/proto"
	"time"
)

func main() {
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.Dial("10.1.3.105:56923", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := service.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	res, err := client.SayHello(ctx, &service.HelloRequest{Name: "jajarun"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	fmt.Println("reply msg:", res.GetMessage())
}
