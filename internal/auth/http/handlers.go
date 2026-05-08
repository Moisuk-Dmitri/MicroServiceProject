package authhttp

import (
	authservice "micro-blog/internal/auth/service"
	"net/http"
)

type Handler struct {
	service authservice.Service
}

func NewHandler(service authservice.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Register(r.Context()); err != nil {
		http.Error(w,
			"internal server error",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user registered"))
}
