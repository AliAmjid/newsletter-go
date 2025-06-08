package postgres

import (
	"context"
	"database/sql"

	"github.com/AliAmjid/newsletter-go/domain"
)

type PostRepository struct {
	DB *sql.DB
}

func NewPostRepository(db *sql.DB) *PostRepository {
	return &PostRepository{DB: db}
}

func (p *PostRepository) Store(ctx context.Context, post *domain.Post) error {
	_, err := p.DB.ExecContext(ctx,
		"INSERT INTO posts (title, content) VALUES ($1, $2)",
		post.Title, post.Content,
	)
	return err
}

func (p *PostRepository) ListByNewsletterID(ctx context.Context, newsletterID string) ([]domain.Post, error) {
	rows, err := p.DB.QueryContext(ctx,
		"SELECT title, content FROM posts WHERE newsletter_id = $1",
		newsletterID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []domain.Post
	for rows.Next() {
		var post domain.Post
		if err := rows.Scan(&post.Title, &post.Content); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	return posts, rows.Err()
}
