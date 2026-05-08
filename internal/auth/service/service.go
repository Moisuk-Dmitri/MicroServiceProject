package authservice

import "context"

type UserCreatedPublisher interface {
	UserCreatedEvent(ctx context.Context) error
}

type Service interface {
	Register(ctx context.Context) error
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type service struct {
	publisher UserCreatedPublisher
}

func NewService(publisher UserCreatedPublisher) Service {
	return &service{
		publisher: publisher,
	}
}

func (s *service) Register(ctx context.Context) error {
	return s.publisher.UserCreatedEvent(ctx)
}

func (s *service) ValidateToken(ctx context.Context, token string) (bool, error) {
	if token == "valid-token" {
		return true, nil
	}

	return false, nil
}
