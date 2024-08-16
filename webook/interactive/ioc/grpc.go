package ioc

import (
	grpc2 "basic-go/webook/interactive/grpc"
	"basic-go/webook/pkg/grpcx"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func NewGrpcxServer(intrSvc *grpc2.InteractiveServiceServer) *grpcx.Server {
	s := grpc.NewServer()
	intrSvc.Register(s)
	addr := viper.GetString("grpc.server.addr")
	return &grpcx.Server{
		Server: s,
		Addr:   addr,
	}
}
