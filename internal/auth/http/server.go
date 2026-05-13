package authhttp

import (
	"context"
	"log"
	authservice "micro-blog/internal/auth/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	addr   string
	server *http.Server
}

func NewServer(port string, service authservice.Service) *Server {
	handler := NewHandler(service)

	router := chi.NewRouter()
	router.Post("/register", handler.Register)

	if port[0] != ':' {
		port = ":" + port
	}

	return &Server{
		addr: port,
		server: &http.Server{
			Addr:    port,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("auth http server started on %s", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
