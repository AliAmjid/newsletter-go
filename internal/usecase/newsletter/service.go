package newsletter

import (
	"context"
	"errors"

	"newsletter-go/domain"
)

// Service provides newsletter business logic.
type Service struct {
	repo domain.NewsletterRepository
}

var ErrNotFound = errors.New("newsletter not found")

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
	n, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if n == nil {
		return nil, ErrNotFound
	}
	return n, nil
}

func (s *Service) Update(ctx context.Context, n *domain.Newsletter) error {
	if _, err := s.GetByID(ctx, n.ID); err != nil {
		return err
	}
	return s.repo.Update(ctx, n)
}

func (s *Service) Delete(ctx context.Context, id string) error {
	if _, err := s.GetByID(ctx, id); err != nil {
		return err
	}
	return s.repo.Delete(ctx, id)
}
