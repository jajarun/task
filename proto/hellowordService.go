package service

import (
	"context"
	"fmt"
)

type HelloworldService struct {
	UnimplementedGreeterServer
}

func (s *HelloworldService) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
	fmt.Println("get request:", in)
	return &HelloReply{Message: "Hello " + in.Name}, nil
}
