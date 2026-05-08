package blogservice

import "context"

type Service interface {
	CreatePost(ctx context.Context) error
}

type service struct{}

func NewService() Service {
	return &service{}
}

func (s *service) CreatePost(ctx context.Context) error {
	return nil
}
