package newsletter

import (
	"context"

	"newsletter-go/domain"
)

// Service provides newsletter business logic.
type Service struct {
	repo domain.NewsletterRepository
}

func NewService(r domain.NewsletterRepository) *Service {
	return &Service{repo: r}
}

func (s *Service) ListByOwner(ctx context.Context, ownerID string) ([]*domain.Newsletter, error) {
	return s.repo.ListByOwner(ctx, ownerID)
}

func (s *Service) Create(ctx context.Context, n *domain.Newsletter) error {
	return s.repo.Create(ctx, n)
}

func (s *Service) GetByID(ctx context.Context, id string) (*domain.Newsletter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) Update(ctx context.Context, n *domain.Newsletter) error {
	return s.repo.Update(ctx, n)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
