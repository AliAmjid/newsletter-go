package domain

import "context"

type Post struct {
	ID           string `json:"id"`
	NewsletterId string `json:"newsletter_id"`
	Title        string `json:"title"`
	Content      string `json:"content"`
	PublishedAt  string `json:"published_at"`
}

type PostRepository interface {
	Store(ctx context.Context, p *Post) error
	Create(ctx context.Context, p *Post) error
	ListPostsByNewsletter(ctx context.Context, newsletterId string) ([]*Post, error)
}
