package authhttp

import (
	"log"
	authservice "micro-blog/internal/auth/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	addr   string
	server *http.Server
}

func NewServer(addr string, service authservice.Service) *Server {
	handler := NewHandler(service)

	router := chi.NewRouter()
	router.Post("/register", handler.Register)

	return &Server{
		addr: addr,
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("auth http server started on %s", s.addr)
	return s.server.ListenAndServe()
}
