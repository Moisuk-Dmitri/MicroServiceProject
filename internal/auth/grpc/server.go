package authgrpc

import (
	"context"
	"log"
	authservice "micro-blog/internal/auth/service"
	authpb "micro-blog/proto/auth"
	"net"

	"google.golang.org/grpc"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer

	addr       string
	service    authservice.Service
	grpcServer *grpc.Server
}

func NewServer(port string, service authservice.Service) *Server {
	if port[0] != ':' {
		port = ":" + port
	}

	return &Server{
		addr:       port,
		service:    service,
		grpcServer: grpc.NewServer(),
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	authpb.RegisterAuthServiceServer(s.grpcServer, s)

	log.Printf("auth grpc server started on %s", s.addr)

	return s.grpcServer.Serve(listener)
}

func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}

func (s *Server) ValidateToken(
	ctx context.Context,
	req *authpb.ValidateTokenRequest,
) (*authpb.ValidateTokenResponse, error) {
	userID, err := s.service.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	return &authpb.ValidateTokenResponse{
		Valid:  userID != "",
		UserId: userID,
	}, nil
}
