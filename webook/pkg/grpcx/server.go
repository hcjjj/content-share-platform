package grpcx

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	*grpc.Server
	Addr string
}

func (s *Server) Serve() error {
	l, err := net.Listen("tcp", s.Addr)
	fmt.Printf("grpc server listen on %s ...", s.Addr)
	if err != nil {
		return err
	}
	// 这边会阻塞，类似与 gin.Run
	return s.Server.Serve(l)
}
