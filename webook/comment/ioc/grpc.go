package ioc

import (
	"basic-go/webook/pkg/grpcx"
	"basic-go/webook/pkg/logger"

	grpc2 "basic-go/webook/comment/grpc"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func InitGrpcxServer(comment *grpc2.CommentServiceServer, l logger.LoggerV1) *grpcx.Server {
	type Config struct {
		EtcdAddr string `yaml:"etcdAddr"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc", &cfg)
	if err != nil {
		panic(err)
	}
	server := grpc.NewServer()
	comment.Register(server)
	return &grpcx.Server{
		Server:   server,
		EtcdAddr: cfg.EtcdAddr,
		Port:     cfg.Port,
		Name:     cfg.Name,
		L:        l,
	}
}
