package bloghttp

import (
	"context"
	"log"
	"net/http"

	blogservice "micro-blog/internal/blog/service"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	addr   string
	server *http.Server
}

func NewServer(
	addr string,
	service blogservice.Service,
	authenticator Authenticator,
) *Server {
	handler := NewHandler(service)

	router := chi.NewRouter()
	router.With(AuthMiddleware(authenticator)).Post("/posts", handler.CreatePost)
	router.With(AuthMiddleware(authenticator)).Get("/posts", handler.GetPosts)

	return &Server{
		addr: addr,
		server: &http.Server{
			Addr:    addr,
			Handler: router,
		},
	}
}

func (s *Server) Run() error {
	log.Printf("blog http server started on %s", s.addr)
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
