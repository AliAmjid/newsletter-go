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

func (p *PostRepository) ListPostsByNewsletter(ctx context.Context, newsletterId string) ([]*domain.Post, error) {
	rows, err := p.DB.QueryContext(ctx,
		"SELECT id, newsletter_id, title, content, published_at FROM post WHERE newsletter_id = $1",
		newsletterId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*domain.Post
	for rows.Next() {
		post := &domain.Post{}
		if err := rows.Scan(&post.ID, &post.NewsletterId, &post.Title, &post.Content, &post.PublishedAt); err != nil {
			return []*domain.Post{}, err
		}
		posts = append(posts, post)
	}

	// Check for any errors encountered during iteration
	if err := rows.Err(); err != nil {
		return []*domain.Post{}, err
	}

	// If no posts were found, return an empty slice
	if posts == nil {
		posts = []*domain.Post{}
	}

	return posts, nil
}
