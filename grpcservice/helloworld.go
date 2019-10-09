package grpcservice

import (
	"context"
	"github.com/myadamtest/adam/grpcservice/pb/adam"
)

type GreeterImpl struct{}

func (this *GreeterImpl) SayHello(ctx context.Context, in *adam.HelloRequest) (*adam.HelloReply, error) {
	panic("to implement")
	return nil, nil
}

func (this *GreeterImpl) SayHello2(st adam.Greeter_SayHello2Server) error {
	panic("to implement")
	return nil
}

func (this *GreeterImpl) SayHello3(st adam.Greeter_SayHello3Server) error {
	panic("to implement")
	return nil
}

func (this *GreeterImpl) SayHello4(in *adam.HelloRequest, st adam.Greeter_SayHello4Server) error {
	panic("to implement")
	return nil
}
