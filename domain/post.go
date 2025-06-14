package domain

import (
	"context"
	"time"
)

type Post struct {
	ID           string     `json:"id"`
	NewsletterId string     `json:"newsletterId"`
	Title        string     `json:"title"`
	Content      string     `json:"content"`
	PublishedAt  *time.Time `json:"publishedAt,omitempty"`
}

type PostRepository interface {
	Create(ctx context.Context, p *Post) error
	GetByID(ctx context.Context, id string) (*Post, error)
	Publish(ctx context.Context, id string) (*Post, error)
	ListByNewsletter(ctx context.Context, newsletterId, cursor string, limit int) ([]*Post, error)
}
