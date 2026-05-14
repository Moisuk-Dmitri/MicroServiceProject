package authservice

import (
	"context"
	"fmt"
	authmodel "micro-blog/internal/auth/model"
	autherrors "micro-blog/internal/auth/model/errors"
	"micro-blog/internal/events"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const accessTokenSecret = "dev-secret"

type UserRepository interface {
	Create(ctx context.Context, user authmodel.User) error
	GetByEmail(ctx context.Context, email string) (authmodel.User, error)
}

type UserCreatedPublisher interface {
	UserCreatedEvent(ctx context.Context, event events.UserCreatedEvent) error
}

type Service interface {
	Register(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
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

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.users.GetByEmail(ctx, email)
	if err != nil {
		return "", autherrors.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", autherrors.ErrInvalidCredentials
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"exp":     time.Now().Add(15 * time.Minute).Unix(),
	}).SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return token, nil
}

func (s *service) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	token, err := jwt.Parse(tokenString, s.keyFunc)
	if err != nil {
		return "", nil
	}
	if !token.Valid {
		return "", nil
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", nil
	}
	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", nil
	}

	return userID, nil
}

func (s *service) keyFunc(token *jwt.Token) (interface{}, error) {
	return []byte(accessTokenSecret), nil
}
