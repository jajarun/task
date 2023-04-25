package serviceMicro

import (
	"context"
	"fmt"
)

type HelloworldService struct {
	GreeterHandler
}

//func (s *HelloworldService) SayHello(ctx context.Context, in *HelloRequest) (*HelloReply, error) {
//	fmt.Println("get request:", in)
//	return &HelloReply{Message: "Hello " + in.Name}, nil
//}

func (s *HelloworldService) SayHello(ctx context.Context, in *HelloRequest, out *HelloReply) error {
	fmt.Println("get request:", in)
	out.Message = "你好啊"
	return nil
	//return h.GreeterHandler.SayHello(ctx, in, out)
}
