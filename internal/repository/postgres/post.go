package postgres

import (
	"context"
	"database/sql"

	"github.com/AliAmjid/newsletter-go/domain"
)

// PostRepository is a PostgreSQL implementation of domain.PostRepository.
type PostRepository struct {
	DB *sql.DB
}

// NewPostRepository creates a new PostRepository.
func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

// Store inserts a post into the database.
func (p *PostRepository) Store(ctx context.Context, post *domain.Post) error {
	_, err := p.DB.ExecContext(ctx,
		"INSERT INTO posts (title, content) VALUES ($1, $2)",
		post.Title, post.Content,
	)
	return err
}
