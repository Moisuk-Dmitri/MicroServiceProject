package bloghttp

import (
	"encoding/json"
	blogdto "micro-blog/internal/blog/dto"
	blogmodel "micro-blog/internal/blog/model"
	blogservice "micro-blog/internal/blog/service"
	"net/http"
	"time"
)

type Handler interface {
	CreatePost(w http.ResponseWriter, r *http.Request)
	GetPosts(w http.ResponseWriter, r *http.Request)
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
	userID, ok := r.Context().Value(UserIDKey).(string)
	if !ok {
		http.Error(
			w,
			"unathorized",
			http.StatusUnauthorized,
		)
		return
	}

	var post blogmodel.Post
	if err := json.NewDecoder(r.Body).Decode(&post); err != nil {
		http.Error(
			w,
			"request decode error",
			http.StatusBadRequest,
		)
		return
	}

	if err := h.service.CreatePost(r.Context(), userID, post.Title, post.Content); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (h *handler) GetPosts(w http.ResponseWriter, r *http.Request) {
	posts, err := h.service.ListPosts(r.Context())
	if err != nil {
		http.Error(
			w,
			"internal server error",
			http.StatusInternalServerError,
		)
		return
	}

	resp := make([]blogdto.PostResponse, 0, len(posts))
	for _, post := range posts {
		resp = append(resp, blogdto.PostResponse{
			ID:        post.ID,
			AuthorID:  post.AuthorID,
			Title:     post.Title,
			Content:   post.Content,
			CreatedAt: post.CreatedAt.Format(time.RFC3339),
		})
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(
			w,
			"response encode error",
			http.StatusInternalServerError,
		)
		return
	}
}
