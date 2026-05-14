package blogservice

import (
	"context"
	"fmt"
	blogmodel "micro-blog/internal/blog/model"
	"time"

	"github.com/google/uuid"
)

type PostRepository interface {
	Create(ctx context.Context, post blogmodel.Post) error
	List(ctx context.Context) ([]blogmodel.Post, error)
}

type Service interface {
	CreatePost(ctx context.Context, authorID, title, content string) error
	ListPosts(ctx context.Context) ([]blogmodel.Post, error)
}

type service struct {
	posts PostRepository
}

func NewService(postRepository PostRepository) Service {
	return &service{
		posts: postRepository,
	}
}

func (s *service) CreatePost(ctx context.Context, authorID, title, content string) error {
	if title == "" {
		return fmt.Errorf("title is required")
	}

	if content == "" {
		return fmt.Errorf("content is required")
	}

	return s.posts.Create(
		ctx,
		blogmodel.Post{
			ID:        uuid.NewString(),
			AuthorID:  authorID,
			Title:     title,
			Content:   content,
			CreatedAt: time.Now(),
		},
	)
}

func (s *service) ListPosts(ctx context.Context) ([]blogmodel.Post, error) {
	return s.posts.List(ctx)
}
