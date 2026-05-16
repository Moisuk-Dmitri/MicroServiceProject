package blogrepository

import (
	"context"
	"database/sql"
	blogmodel "micro-blog/internal/blog/model"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	_ "github.com/lib/pq"
)

func TestRepository_CreateAndList_Success(t *testing.T) {
	db, err := sql.Open("postgres", "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	require.NoError(t, err)
	defer db.Close()

	ctx := context.Background()
	userID := uuid.NewString()
	_, err = db.ExecContext(ctx, `
		INSERT INTO users (id, email, password_hash, created_at)
		VALUES ($1, $2, $3, NOW())
	`, userID, "test@example.com", "hash")

	repo := NewRepository(db)

	post := blogmodel.Post{
		ID:        uuid.NewString(),
		AuthorID:  userID,
		Title:     "title",
		Content:   "content",
		CreatedAt: time.Now(),
	}

	err = repo.Create(ctx, post)
	require.NoError(t, err)

	posts, err := repo.List(ctx)
	require.NoError(t, err)

	var found bool
	for _, got := range posts {
		if got.ID == post.ID {
			found = true

			require.Equal(t, post.Title, got.Title)
			require.Equal(t, post.Content, got.Content)
			require.Equal(t, post.AuthorID, got.AuthorID)

			break
		}
	}

	require.True(t, found, "created post was not found in list")

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DELETE FROM posts WHERE author_id = $1", userID)
		_, _ = db.ExecContext(ctx, "DELETE FROM users WHERE id = $1", userID)
	})
}
