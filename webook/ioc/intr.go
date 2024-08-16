package ioc

import (
	intrv1 "basic-go/webook/api/proto/gen/intr/v1"
	"basic-go/webook/interactive/service"
	"basic-go/webook/internal/client"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitIntrClient(svc service.InteractiveService) intrv1.InteractiveServiceClient {
	type Config struct {
		Addr      string `yaml:"addr"`
		Secure    bool
		Threshold int32
	}
	var cfg Config
	err := viper.UnmarshalKey("grpc.client.intr", &cfg)
	if err != nil {
		panic(err)
	}
	var opts []grpc.DialOption
	if !cfg.Secure {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	cc, err := grpc.Dial(cfg.Addr, opts...)
	if err != nil {
		panic(err)
	}
	remote := intrv1.NewInteractiveServiceClient(cc)
	local := client.NewLocalInteractiveServiceAdapter(svc)
	res := client.NewInteractiveClient(remote, local)
	viper.OnConfigChange(func(in fsnotify.Event) {
		cfg = Config{}
		err := viper.UnmarshalKey("grpc.client.intr", &cfg)
		if err != nil {
			// 这边做不了什么
			panic(err)
		}
		res.UpdateThreshold(cfg.Threshold)
	})
	return res
}
