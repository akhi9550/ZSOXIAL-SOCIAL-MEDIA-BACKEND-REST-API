package server

import (
	"fmt"
	"net"

	"github.com/akhi9550/notification-svc/pkg/config"
	pb "github.com/akhi9550/notification-svc/pkg/pb/notification"

	"google.golang.org/grpc"
)

type Server struct {
	server   *grpc.Server
	listener net.Listener
}

func NewGRPCServer(cfg config.Config, server pb.NotificationServiceServer) (*Server, error) {
	lis, err := net.Listen("tcp", cfg.Port)
	if err != nil {
		return nil, err
	}
	newServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(newServer, server)

	return &Server{
		server:   newServer,
		listener: lis,
	}, nil
}

func (c *Server) Start() error {
	fmt.Println("grpc server listening on port :50054")
	return c.server.Serve(c.listener)
}
