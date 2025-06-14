package post

import (
	"context"

	"newsletter-go/domain"
)

// Service handles post business logic.
type Service struct {
	repo           domain.PostRepository
	newsletterRepo domain.NewsletterRepository
}

// NewService creates a Service instance.
func NewService(r domain.PostRepository, nR domain.NewsletterRepository) *Service {
	return &Service{
		repo:           r,
		newsletterRepo: nR,
	}
}

// Save persists a post via the repository.
func (s *Service) Save(ctx context.Context, p *domain.Post) error {
	return s.repo.Create(ctx, p)
}

func (s *Service) IsNewsletterOwner(ctx context.Context, newsletterId, userId string) (bool, error) {
	return s.newsletterRepo.IsOwner(ctx, newsletterId, userId)
}
