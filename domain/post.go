package domain

import "context"

type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type PostRepository interface {
	Store(ctx context.Context, p *Post) error
}
