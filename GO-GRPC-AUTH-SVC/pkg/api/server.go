package server

import (
	"fmt"
	"net"

	"github.com/akhi9550/auth-svc/pkg/config"
	pb "github.com/akhi9550/auth-svc/pkg/pb/auth"

	"google.golang.org/grpc"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGRPCServer(cfg config.Config, server pb.AuthServiceServer) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return nil, err
	}
	newServer := grpc.NewServer()
	pb.RegisterAuthServiceServer(newServer, server)

	return &Server{
		server:   newServer,
		listener: lis,
	}, nil
}

func (c *Server) Start() error {
	fmt.Println("grpc server listening on port :50051")
	return c.server.Serve(c.listener)
}
