package bloghttp

import (
	blogservice "micro-blog/internal/blog/service"
	"net/http"
)

type Handler interface {
	CreatePost(http.ResponseWriter, *http.Request)
}

type handler struct {
	service blogservice.Service
}

func NewHandler(svc blogservice.Service) Handler {
	return &handler{
		service: svc,
	}
}

func (h *handler) CreatePost(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("posts list"))
}
