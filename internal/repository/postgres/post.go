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

func (p *PostRepository) Create(ctx context.Context, post *domain.Post) error {
	return p.DB.QueryRowContext(ctx,
		"INSERT INTO post (newsletter_id, title, content, published_at) VALUES ($1, $2, $3, $4) RETURNING id, published_at",
		post.NewsletterId, post.Title, post.Content, post.PublishedAt,
	).Scan(&post.ID, &post.PublishedAt)
}

func (p *PostRepository) GetByID(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post
	err := p.DB.QueryRowContext(ctx,
		`SELECT id, newsletter_id, title, content, published_at FROM post WHERE id = $1`,
		id,
	).Scan(&post.ID, &post.NewsletterId, &post.Title, &post.Content, &post.PublishedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (p *PostRepository) Publish(ctx context.Context, id string) (*domain.Post, error) {
	var post domain.Post
	err := p.DB.QueryRowContext(ctx,
		`UPDATE post SET published_at = NOW() WHERE id = $1 RETURNING id, newsletter_id, title, content, published_at`,
		id,
	).Scan(&post.ID, &post.NewsletterId, &post.Title, &post.Content, &post.PublishedAt)
	if err != nil {
		return nil, err
	}
	return &post, nil
}
func (p *PostRepository) ListByNewsletter(ctx context.Context, newsletterId, cursor string, limit int) ([]*domain.Post, error) {
	var rows *sql.Rows
	var err error
	if cursor == "" {
		rows, err = p.DB.QueryContext(ctx,
			`SELECT id, newsletter_id, title, content, published_at
                         FROM post
                         WHERE newsletter_id = $1 AND published_at IS NOT NULL
                         ORDER BY published_at DESC
                         LIMIT $2`,
			newsletterId, limit,
		)
	} else {
		rows, err = p.DB.QueryContext(ctx,
			`SELECT id, newsletter_id, title, content, published_at
                         FROM post
                         WHERE newsletter_id = $1 AND published_at IS NOT NULL AND published_at < $2
                         ORDER BY published_at DESC
                         LIMIT $3`,
			newsletterId, cursor, limit,
		)
	}
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
