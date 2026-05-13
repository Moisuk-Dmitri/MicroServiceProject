package authservice

import (
	"context"
	"errors"
	"fmt"
	authmodel "micro-blog/internal/auth/model"
	autherrors "micro-blog/internal/auth/model/errors"
	"micro-blog/internal/events"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(ctx context.Context, user authmodel.User) error
	GetByEmail(ctx context.Context, email string) (authmodel.User, error)
}

type UserCreatedPublisher interface {
	UserCreatedEvent(ctx context.Context, event events.UserCreatedEvent) error
}

type Service interface {
	Register(ctx context.Context, email, password string) error
	ValidateToken(ctx context.Context, token string) (bool, error)
}

type service struct {
	users     UserRepository
	publisher UserCreatedPublisher
}

func NewService(publisher UserCreatedPublisher, userRepository UserRepository) Service {
	return &service{
		users:     userRepository,
		publisher: publisher,
	}
}

func (s *service) Register(ctx context.Context, email, password string) error {
	_, err := s.users.GetByEmail(ctx, email)
	if err == nil {
		return autherrors.ErrUserAlreadyExists
	}
	if !errors.Is(err, autherrors.ErrUserNotFound) {
		return fmt.Errorf("register user: %w", err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}
	user := authmodel.User{
		ID:           uuid.NewString(),
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
	}

	if err := s.users.Create(ctx, user); err != nil {
		return fmt.Errorf("register user: %w", err)
	}

	return s.publisher.UserCreatedEvent(ctx, events.UserCreatedEvent{
		ID:    user.ID,
		Email: user.Email,
	})
}

func (s *service) ValidateToken(ctx context.Context, token string) (bool, error) {
	return token == "valid-token", nil
}
