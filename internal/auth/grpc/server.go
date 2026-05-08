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

	addr    string
	service authservice.Service
}

func NewServer(addr string, service authservice.Service) *Server {
	return &Server{
		addr:    addr,
		service: service,
	}
}

func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()
	authpb.RegisterAuthServiceServer(grpcServer, s)

	log.Printf("auth grpc server started on %s", s.addr)

	return grpcServer.Serve(listener)
}

func (s *Server) ValidateToken(
	ctx context.Context,
	req *authpb.ValidateTokenRequest,
) (*authpb.ValidateTokenResponse, error) {
	valid, err := s.service.ValidateToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}

	return &authpb.ValidateTokenResponse{
		Valid: valid,
	}, nil
}
