package blogrepository

import (
	"context"
	"database/sql"
	"fmt"
	blogmodel "micro-blog/internal/blog/model"
)

type PostRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *PostRepository {
	return &PostRepository{
		db: db,
	}
}

func (r *PostRepository) Create(ctx context.Context, post blogmodel.Post) error {
	query := `
		INSERT INTO posts (id, author_id, title, content, created_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		post.ID,
		post.AuthorID,
		post.Title,
		post.Content,
		post.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("postgre create post: %w", err)
	}

	return nil
}

func (r *PostRepository) List(ctx context.Context) ([]blogmodel.Post, error) {
	query := `
		SELECT id, author_id, title, content, created_at
		FROM posts
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
	)
	if err != nil {
		return nil, fmt.Errorf("list posts: %w", err)
	}
	defer rows.Close()

	posts := make([]blogmodel.Post, 0)
	for rows.Next() {
		var post blogmodel.Post

		if err := rows.Scan(
			&post.ID,
			&post.AuthorID,
			&post.Title,
			&post.Content,
			&post.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan post: %w", err)
		}

		posts = append(posts, post)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate posts: %w", err)
	}

	return posts, nil
}
