package authhttp

import (
	"encoding/json"
	"errors"
	authdto "micro-blog/internal/auth/dto"
	autherrors "micro-blog/internal/auth/model/errors"
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
	registerRequest := &authdto.RegisterRequest{}
	err := json.NewDecoder(r.Body).Decode(registerRequest)
	if err != nil {
		http.Error(
			w,
			"invalid request body",
			http.StatusBadRequest,
		)
		return
	}
	if registerRequest.Email == "" ||
		registerRequest.Password == "" {
		http.Error(
			w,
			"empty email or password",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.Register(
		r.Context(),
		registerRequest.Email,
		registerRequest.Password); err != nil {
		if errors.Is(err, autherrors.ErrUserAlreadyExists) {
			http.Error(
				w,
				"user already exists",
				http.StatusConflict)
			return
		}

		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("user registered"))
}
