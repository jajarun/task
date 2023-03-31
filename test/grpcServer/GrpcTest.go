package main

import (
	"fmt"
	"google.golang.org/grpc"
	"net"
	service "task/proto"
)

func main() {
	lis, err := net.Listen("tcp", ":8972")
	if err != nil {
		fmt.Printf("failed to listen: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	service.RegisterGreeterServer(grpcServer, &service.HelloworldService{})
	err = grpcServer.Serve(lis)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
}
