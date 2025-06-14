package service

import (
	"context"
	"newsletter-go/domain"
)

type NewsletterService struct {
	repo domain.NewsletterRepository
}

func NewNewsletterService(repo domain.NewsletterRepository) *NewsletterService {
	return &NewsletterService{repo: repo}
}

func (s *NewsletterService) Create(ctx context.Context, n *domain.Newsletter) error {
	return s.repo.Create(ctx, n)
}

func (s *NewsletterService) GetByID(ctx context.Context, id string) (*domain.Newsletter, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *NewsletterService) ListByOwner(ctx context.Context, ownerID string) ([]*domain.Newsletter, error) {
	return s.repo.ListByOwner(ctx, ownerID)
}

func (s *NewsletterService) Update(ctx context.Context, n *domain.Newsletter) error {
	return s.repo.Update(ctx, n)
}

func (s *NewsletterService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *NewsletterService) IsOwner(ctx context.Context, newsletterId, userId string) (bool, error) {
	return s.repo.IsOwner(ctx, newsletterId, userId)
}
