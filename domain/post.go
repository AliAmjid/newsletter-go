package domain

import "context"

// Post represents a newsletter post.
type Post struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// PostRepository defines the persistence functions for posts.
type PostRepository interface {
	Store(ctx context.Context, p *Post) error
}
