package postgres

import (
	"context"
	"database/sql"

	"newsletter-go/domain"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (p *PostRepository) Store(ctx context.Context, post *domain.Post) error {
	_, err := p.DB.ExecContext(ctx,
		"INSERT INTO post (title, content) VALUES ($1, $2)",
		post.Title, post.Content,
	)
	return err
}

func (p *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return p.DB.QueryRowContext(ctx,
		"INSERT INTO post (newsletter_id, title, content) VALUES ($1, $2, $3) RETURNING id, published_at",
		post.NewsletterId, post.Title, post.Content,
	).Scan(&post.ID, &post.PublishedAt)
}
