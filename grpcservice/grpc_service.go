package grpcservice

import (
	"fmt"
	"github.com/myadamtest/adam/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

func StartGrpc() error {
	s := grpc.NewServer()

	//...

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", config.GetConfig().RpcPort))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	reflection.Register(s)

	err = s.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	return nil
}
