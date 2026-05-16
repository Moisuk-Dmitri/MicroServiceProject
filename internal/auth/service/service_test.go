package authservice

import (
	"context"
	"errors"
	authmodel "micro-blog/internal/auth/model"
	autherrors "micro-blog/internal/auth/model/errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	user authmodel.User
	err  error
}

func (r *MockUserRepository) GetByEmail(ctx context.Context, email string) (authmodel.User, error) {
	if r.err != nil {
		return authmodel.User{}, r.err
	}

	return r.user, nil
}

func (r *MockUserRepository) Create(ctx context.Context, user authmodel.User) error {
	return nil
}

func TestService_Login_Success(t *testing.T) {
	password := "password-test"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)
	repo := &MockUserRepository{
		user: authmodel.User{
			ID:           uuid.NewString(),
			Email:        "email-test",
			PasswordHash: string(hash),
			CreatedAt:    time.Now(),
		},
	}
	svc := NewService(nil, repo)

	ctx := context.Background()

	token, err := svc.Login(ctx, repo.user.Email, password)
	require.NoError(t, err)

	require.NotNil(t, token)
}

func TestService_Login_UserNotFound(t *testing.T) {
	repo := &MockUserRepository{
		err: errors.New("user not found"),
	}

	svc := NewService(nil, repo)

	token, err := svc.Login(context.Background(), "wrong@email.com", "password")

	require.Error(t, err)
	require.Equal(t, autherrors.ErrInvalidCredentials, err)
	require.Empty(t, token)
}
