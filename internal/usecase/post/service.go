package post

import (
	"context"

	"newsletter-go/domain"
)

// Service handles post business logic.
type Service struct {
	repo domain.PostRepository
}

// NewService creates a Service instance.
func NewService(r domain.PostRepository) *Service {
	return &Service{repo: r}
}

// Save persists a post via the repository.
func (s *Service) Save(ctx context.Context, p *domain.Post) error {
	return s.repo.Store(ctx, p)
}
